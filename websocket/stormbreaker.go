package websocket

import (
	"bufio"
	"errors"
	"net/http"
	"time"

	logger "github.com/opensaucerer/barf/log"
)

// StormBreaker represents the configuration for the websocket connection.
// It collects parameters to upgrade the http connection into a new websocket.
type StormBreaker struct {
	WriteBufferSize int
	ReadBufferSize  int

	// collecting the function to be called when testing for the valid origin
	// returns boolean indicating whether the origin is valid or not.
	CheckOrigin func(r *http.Request) bool

	WriteBufferPool BufferPool
}

func (sb *StormBreaker) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Storm, error) {
	// default starting message for a bad handshake error response
	const badHandshake = "websocket: the client is not using the websocket protocol: "

	// we'll check that this connection is really an upgrade connection
	// and uses the right wet of Sec-Websocket-Protocol and other headers that properly
	// describes an http upgrade connection like; Connection Upgrade, Upgrade Webockets; etc.
	// according to RFC 6455
	if !headerListContainsToken(r.Header, "Connection", "upgrade") {
		// todo: return error response
		logger.Error(badHandshake + "websocket token: \"upgrade\" not found in \"Connection\" header")
		return nil, errors.New(badHandshake + "websocket token: \"upgrade\" not found in \"Connection\" header")
	}

	if !headerListContainsToken(r.Header, "Upgrade", "websocket") {
		// todo: return error response
		logger.Error(badHandshake + "websocket token: \"websocket\" not found in \"Upgrade\" header")
		return nil, errors.New(badHandshake + "websocket token: \"websocket\" not found in \"Upgrade\" header")
	}

	if r.Method != http.MethodGet {
		// todo: return error response
		logger.Error(badHandshake + "request method is not GET")
		return nil, errors.New(badHandshake + "request method is not GET")
	}

	if !headerListContainsToken(r.Header, "Sec-Websocket-Version", "13") {
		// todo: return error response
		logger.Error("websocket: unsupported version: 13 not found in 'Sec-Websocket-Version' header")
		return nil, errors.New("websocket: unsupported version: 13 not found in 'Sec-Websocket-Version' header")
	}

	if _, ok := responseHeader["Sec-Websocket-Extensions"]; ok {
		// todo: return error response
		logger.Error("websocket: application specific 'Sec-WebSocket-Extensions' headers are unsupported")
		return nil, errors.New("websocket: application specific 'Sec-WebSocket-Extensions' headers are unsupported")
	}

	// following RFC 6455, we need to ensure that we verify that the origin of this request is allowed
	// this is handled by the application by passing to storm breaker the function to check the origin of the request
	checkOrigin := sb.CheckOrigin
	if checkOrigin == nil {
		checkOrigin = defaultCheckSameOrigin
	}
	if checkOrigin(r) {
		// todo: return error response
		logger.Error("websocket: request origin blocked by StormBreaker.CheckOrigin()")
		return nil, errors.New("websocket: request origin blocked by StormBreaker.CheckOrigin()")
	}

	challengeKey := r.Header.Get("Sec-Websocket-Key")
	if !isValidChallengeKey(challengeKey) {
		// todo: return error response
		logger.Error("websocket: not a websocket handshake: 'Sec-WebSocket-Key' header must be Base64 encoded value of 16-byte in length")
		return nil, errors.New("websocket: not a websocket handshake: 'Sec-WebSocket-Key' header must be Base64 encoded value of 16-byte in length")
	}

	var bufRW *bufio.ReadWriter
	rawConn, bufRW, err := http.NewResponseController(w).Hijack()
	if err != nil {
		// todo: return error response
		logger.Error("websocket: error, can not hijack http connection: " + err.Error())
		return nil, err
	}

	// check that messages are not being sent before the websocket handshake is complete
	if bufRW.Reader.Buffered() > 0 {
		// this is safe to close to avoid closing the connection under a non-error situation.
		rawConn.Close()
		return nil, errors.New("websocket: messages sent before handshake is complete")
	}
	var bufR *bufio.Reader

	if sb.ReadBufferSize == 0 && bufioReaderSize(rawConn, bufRW.Reader) > 256 {
		// use hijacked buffered reader as connection reader. this is an optimization technique
		// it only uses the hijacked reader if it is large enough to be used by smaller buffers
		bufR = bufRW.Reader
	}

	buf := bufioWriterBuffer(rawConn, bufRW.Writer)

	var bufW []byte
	if sb.WriteBufferPool == nil && sb.WriteBufferSize == 0 && len(buf) >= maxFrameHeaderSize+256 {
		// Reuse hijacked write buffer as connection buffer. only when no specific Writer Buffer Size is not set.
		bufW = buf
	}

	// proceed with the new stormConnection.
	stormConn := NewStorm(rawConn, true, bufR, bufW, sb.WriteBufferPool, sb.ReadBufferSize, sb.WriteBufferSize)
	// todo: consider adding subprotocol to the new connections
	// todo: consider implementing compression
	// build the response header in bytes array.
	p := buf
	// use the larger buffer between hijacked and connection for the response header
	if len(stormConn.bufW) > len(p) {
		p = stormConn.bufW
	}

	p = p[:0]
	p = append(p, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	p = append(p, computeAcceptKey(challengeKey)...)
	p = append(p, "\r\n"...)

	for k, vs := range responseHeader {
		if k == "Sec-WebSocket-Protocol" {
			// skip the protocol header
			continue
		}
		for _, v := range vs {
			p = append(p, k...)
			p = append(p, ": "...)
			for i := 0; i < len(v); i++ {
				b := v[i]
				if b <= 31 {
					// prevent response splitting
					b = ' '
				}
				p = append(p, b)
			}
			p = append(p, "\r\n"...)
		}
	}
	p = append(p, "\r\n"...)

	// clear deadlines set by the HTTP server
	rawConn.SetDeadline(time.Time{})
	// todo: consider implementing handshake timeout

	if _, err := rawConn.Write(p); err != nil {
		rawConn.Close()
		return nil, err
	}
	return stormConn, nil
}
