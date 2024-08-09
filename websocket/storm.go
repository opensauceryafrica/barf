package websocket

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// BufferPool represents a pool of buffers. The *sync.Pool type satisfies this
// interface.  The type of the value stored in a pool is not specified.
type BufferPool interface {
	// Get gets a value from the pool or returns nil if the pool is empty.
	Get() interface{}
	// Put adds a value to the pool.
	Put(interface{})
}

// writePoolData is the type added to the write buffer pool. This wrapper is
// used to prevent applications from peeking at and depending on the values
// added to the pool.
type writePoolData struct{ buf []byte }

type messageWriter struct {
	s         *Storm
	compress  bool // whether next call to flushFrame should set RSV1
	pos       int  // end of data in bufW.
	frameType int  // type of the current frame.
	err       error
}

func (w *messageWriter) endMessage(err error) error {
	if w.err != nil {
		return err
	}
	s := w.s
	w.err = err
	s.writer = nil
	if s.writePool != nil {
		s.writePool.Put(writePoolData{buf: s.bufW})
		s.bufW = nil
	}
	return err
}

// flushFrame writes buffered data and extra as a frame to the network. The
// final argument indicates that this is the last frame in the message.
func (w *messageWriter) flushFrame(final bool, extra []byte) error {
	s := w.s
	length := w.pos - maxFrameHeaderSize + len(extra)

	// Check for invalid control frames.
	if isControl(w.frameType) &&
		(!final || length > maxControlFramePayloadSize) {
		return w.endMessage(errInvalidControlFrame)
	}

	b0 := byte(w.frameType)
	if final {
		b0 |= finalBit
	}
	if w.compress {
		b0 |= rsv1Bit
	}
	w.compress = false

	b1 := byte(0)
	if !s.isServer {
		b1 |= maskBit
	}

	// Assume that the frame starts at beginning of s.bufW.
	framePos := 0
	if s.isServer {
		// Adjust up if mask not included in the header.
		framePos = 4
	}

	switch {
	case length >= 65536:
		s.bufW[framePos] = b0
		s.bufW[framePos+1] = b1 | 127
		binary.BigEndian.PutUint64(s.bufW[framePos+2:], uint64(length))
	case length > 125:
		framePos += 6
		s.bufW[framePos] = b0
		s.bufW[framePos+1] = b1 | 126
		binary.BigEndian.PutUint16(s.bufW[framePos+2:], uint16(length))
	default:
		framePos += 8
		s.bufW[framePos] = b0
		s.bufW[framePos+1] = b1 | byte(length)
	}

	if !s.isServer {
		key := newMaskKey()
		copy(s.bufW[maxFrameHeaderSize-4:], key[:])
		maskBytes(key, 0, s.bufW[maxFrameHeaderSize:w.pos])
		if len(extra) > 0 {
			return w.endMessage(s.writeFatal(errors.New("websocket: internal error, extra used in client mode")))
		}
	}

	// Write the buffers to the connection with best-effort detection of
	// concurrent writes. See the concurrency section in the package
	// documentation for more info.

	if s.isWriting {
		panic("concurrent write to websocket connection")
	}
	s.isWriting = true

	err := s.write(w.frameType, s.writeDeadline, s.bufW[framePos:w.pos], extra)

	if !s.isWriting {
		panic("concurrent write to websocket connection")
	}
	s.isWriting = false

	if err != nil {
		return w.endMessage(err)
	}

	if final {
		w.endMessage(errWriteClosed)
		return nil
	}

	// Setup for next frame.
	w.pos = maxFrameHeaderSize
	w.frameType = continuationFrame
	return nil
}

func (w *messageWriter) ncopy(max int) (int, error) {
	n := len(w.s.bufW) - w.pos
	if n <= 0 {
		if err := w.flushFrame(false, nil); err != nil {
			return 0, err
		}
		n = len(w.s.bufW) - w.pos
	}
	if n > max {
		n = max
	}
	return n, nil
}

func (w *messageWriter) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	if len(p) > 2*len(w.s.bufW) && w.s.isServer {
		// Don't buffer large messages.
		err := w.flushFrame(false, p)
		if err != nil {
			return 0, err
		}
		return len(p), nil
	}

	nn := len(p)
	for len(p) > 0 {
		n, err := w.ncopy(len(p))
		if err != nil {
			return 0, err
		}
		copy(w.s.bufW[w.pos:], p[:n])
		w.pos += n
		p = p[n:]
	}
	return nn, nil
}

func (w *messageWriter) WriteString(p string) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	nn := len(p)
	for len(p) > 0 {
		n, err := w.ncopy(len(p))
		if err != nil {
			return 0, err
		}
		copy(w.s.bufW[w.pos:], p[:n])
		w.pos += n
		p = p[n:]
	}
	return nn, nil
}

func (w *messageWriter) ReadFrom(r io.Reader) (nn int64, err error) {
	if w.err != nil {
		return 0, w.err
	}
	for {
		if w.pos == len(w.s.bufW) {
			err = w.flushFrame(false, nil)
			if err != nil {
				break
			}
		}
		var n int
		n, err = r.Read(w.s.bufW[w.pos:])
		w.pos += n
		nn += int64(n)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
	}
	return nn, err
}

func (w *messageWriter) Close() error {
	if w.err != nil {
		return w.err
	}
	return w.flushFrame(true, nil)
}

type messageReader struct{ s *Storm }

func (r *messageReader) Read(b []byte) (int, error) {
	s := r.s
	if s.messageReader != r {
		return 0, io.EOF
	}

	for s.readErr == nil {

		if s.readRemaining > 0 {
			if int64(len(b)) > s.readRemaining {
				b = b[:s.readRemaining]
			}
			n, err := s.bufR.Read(b)
			s.readErr = hideTempErr(err)
			if s.isServer {
				s.readMaskPos = maskBytes(s.readMaskKey, s.readMaskPos, b[:n])
			}
			rem := s.readRemaining
			rem -= int64(n)
			s.setReadRemaining(rem)
			if s.readRemaining > 0 && s.readErr == io.EOF {
				s.readErr = errUnexpectedEOF
			}
			return n, s.readErr
		}

		if s.readFinal {
			s.messageReader = nil
			return 0, io.EOF
		}

		frameType, err := s.advanceFrame()
		switch {
		case err != nil:
			s.readErr = hideTempErr(err)
		case frameType == TextMessage || frameType == BinaryMessage:
			s.readErr = errors.New("websocket: internal error, unexpected text or binary in Reader")
		}
	}

	err := s.readErr
	if err == io.EOF && s.messageReader == r {
		err = errUnexpectedEOF
	}
	return 0, err
}

func (r *messageReader) Close() error {
	return nil
}

// Close codes defined in RFC 6455, section 11.7.
const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)

// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

// ErrCloseSent is returned when the application writes a message to the
// connection after sending a close message.
var ErrCloseSent = errors.New("websocket: close sent")

// ErrReadLimit is returned when reading a message that is larger than the
// read limit set for the connection.
var ErrReadLimit = errors.New("websocket: read limit exceeded")

// netError satisfies the net Error interface.
type netError struct {
	msg       string
	temporary bool
	timeout   bool
}

func (e *netError) Error() string   { return e.msg }
func (e *netError) Temporary() bool { return e.temporary }
func (e *netError) Timeout() bool   { return e.timeout }

// CloseError represents a close message.
type CloseError struct {
	// Code is defined in RFC 6455, section 11.7.
	Code int

	// Text is the optional text payload.
	Text string
}

func (e *CloseError) Error() string {
	s := []byte("websocket: close ")
	s = strconv.AppendInt(s, int64(e.Code), 10)
	switch e.Code {
	case CloseNormalClosure:
		s = append(s, " (normal)"...)
	case CloseGoingAway:
		s = append(s, " (going away)"...)
	case CloseProtocolError:
		s = append(s, " (protocol error)"...)
	case CloseUnsupportedData:
		s = append(s, " (unsupported data)"...)
	case CloseNoStatusReceived:
		s = append(s, " (no status)"...)
	case CloseAbnormalClosure:
		s = append(s, " (abnormal closure)"...)
	case CloseInvalidFramePayloadData:
		s = append(s, " (invalid payload data)"...)
	case ClosePolicyViolation:
		s = append(s, " (policy violation)"...)
	case CloseMessageTooBig:
		s = append(s, " (message too big)"...)
	case CloseMandatoryExtension:
		s = append(s, " (mandatory extension missing)"...)
	case CloseInternalServerErr:
		s = append(s, " (internal server error)"...)
	case CloseTLSHandshake:
		s = append(s, " (TLS handshake error)"...)
	}
	if e.Text != "" {
		s = append(s, ": "...)
		s = append(s, e.Text...)
	}
	return string(s)
}

// IsCloseError returns boolean indicating whether the error is a *CloseError
// with one of the specified codes.
func IsCloseError(err error, codes ...int) bool {
	if e, ok := err.(*CloseError); ok {
		for _, code := range codes {
			if e.Code == code {
				return true
			}
		}
	}
	return false
}

// IsUnexpectedCloseError returns boolean indicating whether the error is a
// *CloseError with a code not in the list of expected codes.
func IsUnexpectedCloseError(err error, expectedCodes ...int) bool {
	if e, ok := err.(*CloseError); ok {
		for _, code := range expectedCodes {
			if e.Code == code {
				return false
			}
		}
		return true
	}
	return false
}

var validReceivedCloseCodes = map[int]bool{
	// see http://www.iana.org/assignments/websocket/websocket.xhtml#close-code-number

	CloseNormalClosure:           true,
	CloseGoingAway:               true,
	CloseProtocolError:           true,
	CloseUnsupportedData:         true,
	CloseNoStatusReceived:        false,
	CloseAbnormalClosure:         false,
	CloseInvalidFramePayloadData: true,
	ClosePolicyViolation:         true,
	CloseMessageTooBig:           true,
	CloseMandatoryExtension:      true,
	CloseInternalServerErr:       true,
	CloseServiceRestart:          true,
	CloseTryAgainLater:           true,
	CloseTLSHandshake:            false,
}

func isValidReceivedCloseCode(code int) bool {
	return validReceivedCloseCodes[code] || (code >= 3000 && code <= 4999)
}

var (
	errWriteTimeout        = &netError{msg: "websocket: write timeout", timeout: true, temporary: true}
	errUnexpectedEOF       = &CloseError{Code: CloseAbnormalClosure, Text: io.ErrUnexpectedEOF.Error()}
	errBadWriteOpCode      = errors.New("websocket: bad write message type")
	errWriteClosed         = errors.New("websocket: write closed")
	errInvalidControlFrame = errors.New("websocket: invalid control frame")
)

func newMaskKey() [4]byte {
	n := rand.Uint32()
	return [4]byte{byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 24)}
}

func hideTempErr(err error) error {
	if e, ok := err.(net.Error); ok && e.Temporary() {
		err = &netError{msg: e.Error(), timeout: e.Timeout()}
	}
	return err
}

func isControl(frameType int) bool {
	return frameType == CloseMessage || frameType == PingMessage || frameType == PongMessage
}

func isData(frameType int) bool {
	return frameType == TextMessage || frameType == BinaryMessage
}

// Storm represents the new Websocket Connection created by the StormBreaker;
// holds the connection information from the original http.Connection
type Storm struct {
	// instance of the hijacked raw http connection
	RawConn  net.Conn
	isServer bool

	// writing fields
	bufW          []byte
	mu            chan struct{} // mutex lock for writing
	writePool     BufferPool
	writeBufSize  int
	writeDeadline time.Time
	writer        io.WriteCloser // the current writer returned to the application
	isWriting     bool           // for best-effort concurrent write detection

	writeErrMu sync.Mutex
	writeErr   error

	//reading fields
	bufR         *bufio.Reader
	readLength   int64 // Message size.
	readFinal    bool  // true the current message has more frames.
	readLimit    int64 // Maximum message size.
	readMaskKey  [4]byte
	readMaskPos  int
	reader       io.ReadCloser // the current reader returned to the application
	readErr      error
	readErrCount int
	handlePong   func(string) error
	handlePing   func(string) error
	handleClose  func(int, string) error
	// bytes remaining in current frame.
	// set setReadRemaining to safely update this value and prevent overflow
	readRemaining int64
	messageReader *messageReader // the current low-level reader

}

func NewStorm(rawConn net.Conn, isServer bool, bufR *bufio.Reader, bufW []byte, writeBufferPool BufferPool, readBufferSize, writeBufferSize int) *Storm {
	if bufR == nil {
		if readBufferSize == 0 {
			readBufferSize = defaultReadBufferSize
		} else if readBufferSize < maxControlFramePayloadSize {
			// must be large enough for control frame
			readBufferSize = maxControlFramePayloadSize
		}
		bufR = bufio.NewReaderSize(rawConn, readBufferSize)

	}
	if writeBufferSize <= 0 {
		writeBufferSize = defaultWriteBufferSize
	}
	writeBufferSize += maxFrameHeaderSize

	if bufW == nil && writeBufferPool == nil {
		bufW = make([]byte, writeBufferSize)
	}

	mu := make(chan struct{}, 1)
	mu <- struct{}{}
	storm := &Storm{
		isServer:     isServer,
		RawConn:      rawConn,
		bufW:         bufW,
		mu:           mu,
		bufR:         bufR,
		readFinal:    true,
		writePool:    writeBufferPool,
		writeBufSize: writeBufferSize,
	}

	return storm
}

func (s *Storm) Close() error {
	return s.RawConn.Close()
}

// setReadRemaining tracks the number of bytes remaining on the connection. If n
// overflows, an ErrReadLimit is returned.
func (s *Storm) setReadRemaining(n int64) error {
	if n < 0 {
		return ErrReadLimit
	}

	s.readRemaining = n
	return nil
}

func (s *Storm) read(n int) ([]byte, error) {
	p, err := s.bufR.Peek(n)
	if err == io.EOF {
		err = errUnexpectedEOF
	}
	s.bufR.Discard(len(p))
	return p, err
}

func (s *Storm) advanceFrame() (int, error) {
	// 1. Skip remainder of previous frame.

	if s.readRemaining > 0 {
		if _, err := io.CopyN(io.Discard, s.bufR, s.readRemaining); err != nil {
			return noFrame, err
		}
	}

	// 2. Read and parse first two bytes of frame header.
	// To aid debugging, collect and report all errors in the first two bytes
	// of the header.

	var errors []string

	p, err := s.read(2)
	if err != nil {
		return noFrame, err
	}
	frameType := int(p[0] & 0xf)
	final := p[0]&finalBit != 0
	rsv1 := p[0]&rsv1Bit != 0
	rsv2 := p[0]&rsv2Bit != 0
	rsv3 := p[0]&rsv3Bit != 0
	mask := p[1]&maskBit != 0

	s.setReadRemaining(int64(p[1] & 0x7f))
	if rsv1 {
		errors = append(errors, "RSV1 set")
	}

	if rsv2 {
		errors = append(errors, "RSV2 set")
	}

	if rsv3 {
		errors = append(errors, "RSV3 set")
	}

	switch frameType {
	case CloseMessage, PingMessage, PongMessage:
		if s.readRemaining > maxControlFramePayloadSize {
			errors = append(errors, "len > 125 for control")
		}
		if !final {
			errors = append(errors, "FIN not set on control")
		}
	case TextMessage, BinaryMessage:
		if !s.readFinal {
			errors = append(errors, "data before FIN")
		}
		s.readFinal = final
	case continuationFrame:
		if s.readFinal {
			errors = append(errors, "continuation after FIN")
		}
		s.readFinal = final
	default:
		errors = append(errors, "bad opcode "+strconv.Itoa(frameType))
	}

	if mask != s.isServer {
		errors = append(errors, "bad MASK")
	}

	if len(errors) > 0 {
		return noFrame, s.handleProtocolError(strings.Join(errors, ", "))
	}

	// 3. Read and parse frame length as per
	// https://tools.ietf.org/html/rfc6455#section-5.2
	//
	// The length of the "Payload data", in bytes: if 0-125, that is the payload
	// length.
	// - If 126, the following 2 bytes interpreted as a 16-bit unsigned
	// integer are the payload length.
	// - If 127, the following 8 bytes interpreted as
	// a 64-bit unsigned integer (the most significant bit MUST be 0) are the
	// payload length. Multibyte length quantities are expressed in network byte
	// order.

	switch s.readRemaining {
	case 126:
		p, err := s.read(2)
		if err != nil {
			return noFrame, err
		}

		if err := s.setReadRemaining(int64(binary.BigEndian.Uint16(p))); err != nil {
			return noFrame, err
		}
	case 127:
		p, err := s.read(8)
		if err != nil {
			return noFrame, err
		}

		if err := s.setReadRemaining(int64(binary.BigEndian.Uint64(p))); err != nil {
			return noFrame, err
		}
	}

	// 4. Handle frame masking.

	if mask {
		s.readMaskPos = 0
		p, err := s.read(len(s.readMaskKey))
		if err != nil {
			return noFrame, err
		}
		copy(s.readMaskKey[:], p)
	}

	// 5. For text and binary messages, enforce read limit and return.

	if frameType == continuationFrame || frameType == TextMessage || frameType == BinaryMessage {

		s.readLength += s.readRemaining
		// Don't allow readLength to overflow in the presence of a large readRemaining
		// counter.
		if s.readLength < 0 {
			return noFrame, ErrReadLimit
		}

		if s.readLimit > 0 && s.readLength > s.readLimit {
			// todo: s.WriteControl(CloseMessage, FormatCloseMessage(CloseMessageTooBig, ""), time.Now().Add(writeWait))
			return noFrame, ErrReadLimit
		}

		return frameType, nil
	}

	// 6. Read control frame payload.

	var payload []byte
	if s.readRemaining > 0 {
		payload, err = s.read(int(s.readRemaining))
		s.setReadRemaining(0)
		if err != nil {
			return noFrame, err
		}
		if s.isServer {
			maskBytes(s.readMaskKey, 0, payload)
		}
	}

	// 7. Process control frame payload.

	switch frameType {
	case PongMessage:
		if err := s.handlePong(string(payload)); err != nil {
			return noFrame, err
		}
	case PingMessage:
		if err := s.handlePing(string(payload)); err != nil {
			return noFrame, err
		}
	case CloseMessage:
		closeCode := CloseNoStatusReceived
		closeText := ""
		if len(payload) >= 2 {
			closeCode = int(binary.BigEndian.Uint16(payload))
			if !isValidReceivedCloseCode(closeCode) {
				return noFrame, s.handleProtocolError("bad close code " + strconv.Itoa(closeCode))
			}
			closeText = string(payload[2:])
			if !utf8.ValidString(closeText) {
				return noFrame, s.handleProtocolError("invalid utf8 payload in close frame")
			}
		}
		if err := s.handleClose(closeCode, closeText); err != nil {
			return noFrame, err
		}
		return noFrame, &CloseError{Code: closeCode, Text: closeText}
	}

	return frameType, nil

}

// SetCloseHandler sets the handler for close messages received from the peer.
// The code argument to h is the received close code or CloseNoStatusReceived
// if the close message is empty. The default close handler sends a close
// message back to the peer.
//
// The handler function is called from the NextReader, ReadMessage and message
// reader Read methods. The application must read the connection to process
// close messages as described in the section on Control Messages above.
//
// The connection read methods return a CloseError when a close message is
// received. Most applications should handle close messages as part of their
// normal error handling. Applications should only set a close handler when the
// application must perform some action before sending a close message back to
// the peer.
func (s *Storm) SetCloseHandler(h func(code int, text string) error) {
	if h == nil {
		h = func(code int, text string) error {
			// todo: message := FormatCloseMessage(code, "")
			// todo: s.WriteControl(CloseMessage, message, time.Now().Add(writeWait))
			return nil
		}
	}
	s.handleClose = h
}

func (s *Storm) handleProtocolError(message string) error {
	data := FormatCloseMessage(CloseProtocolError, message)
	if len(data) > maxControlFramePayloadSize {
		data = data[:maxControlFramePayloadSize]
	}
	// s.WriteControl(CloseMessage, data, time.Now().Add(writeWait))
	return errors.New("websocket: " + message)
}

func (s *Storm) NextReader() (messageType int, r io.Reader, err error) {
	// Close previous Reader; only necessarey for decompression
	if s.reader != nil {
		s.reader.Close()
		s.reader = nil
	}

	s.messageReader = nil
	s.readLength = 0

	for s.readErr == nil {
		frameType, err := s.advanceFrame()
		if err != nil {
			s.readErr = hideTempErr(err)
			break
		}

		if frameType == TextMessage || frameType == BinaryMessage {
			s.messageReader = &messageReader{s}
			s.reader = s.messageReader
			// if s.readDecompress {
			// 	s.reader = s.newDecompressionReader(s.reader)
			// }
			return frameType, s.reader, nil
		}
	}

	// Applications that do handle the error returned from this method spin in
	// tight loop on connection failure. To help application developers detect
	// this error, panic on repeated reads to the failed connection.
	s.readErrCount++
	if s.readErrCount >= 1000 {
		panic("repeated read on failed websocket connection")
	}

	return noFrame, nil, s.readErr
}

// ReadMessage is a helper method for getting a reader using NextReader and
// reading from that reader to a buffer.
func (s *Storm) ReadMessage() (messageType int, p []byte, err error) {
	var r io.Reader
	messageType, r, err = s.NextReader()
	if err != nil {
		return messageType, nil, err
	}
	p, err = io.ReadAll(r)
	return messageType, p, err
}

// beginMessage prepares a connection and message writer for a new message.
func (s *Storm) beginMessage(mw *messageWriter, messageType int) error {
	// Close previous writer if not already closed by the application. It's
	// probably better to return an error in this situation, but we cannot
	// change this without breaking existing applications.
	if s.writer != nil {
		s.writer.Close()
		s.writer = nil
	}

	if !isControl(messageType) && !isData(messageType) {
		return errBadWriteOpCode
	}

	s.writeErrMu.Lock()
	err := s.writeErr
	s.writeErrMu.Unlock()
	if err != nil {
		return err
	}

	mw.s = s
	mw.frameType = messageType
	mw.pos = maxFrameHeaderSize

	if s.bufW == nil {
		wpd, ok := s.writePool.Get().(writePoolData)
		if ok {
			s.bufW = wpd.buf
		} else {
			s.bufW = make([]byte, s.writeBufSize)
		}
	}
	return nil
}

// NextWriter returns a writer for the next message to send. The writer's Close
// method flushes the complete message to the network.
//
// There can be at most one open writer on a connection. NextWriter closes the
// previous writer if the application has not already done so.
//
// All message types (TextMessage, BinaryMessage, CloseMessage, PingMessage and
// PongMessage) are supported.
func (s *Storm) NextWriter(messageType int) (io.WriteCloser, error) {
	var mw messageWriter
	if err := s.beginMessage(&mw, messageType); err != nil {
		return nil, err
	}
	s.writer = &mw
	return s.writer, nil
}

// WriteMessage is a helper method for getting a writer using NextWriter,
// writing the message and closing the writer.
func (s *Storm) WriteMessage(messageType int, data []byte) error {

	// if s.isServer && (s.newCompressionWriter == nil || !s.enableWriteCompression) {
	if s.isServer {
		// Fast path with no allocations and single frame.

		var mw messageWriter
		if err := s.beginMessage(&mw, messageType); err != nil {
			return err
		}
		n := copy(s.bufW[mw.pos:], data)
		mw.pos += n
		data = data[n:]
		return mw.flushFrame(true, data)
	}

	w, err := s.NextWriter(messageType)
	if err != nil {
		return err
	}
	if _, err = w.Write(data); err != nil {
		return err
	}
	return w.Close()
}

// FormatCloseMessage formats closeCode and text as a WebSocket close message.
// An empty message is returned for code CloseNoStatusReceived.
func FormatCloseMessage(closeCode int, text string) []byte {
	if closeCode == CloseNoStatusReceived {
		// Return empty message because it's illegal to send
		// CloseNoStatusReceived. Return non-nil value in case application
		// checks for nil.
		return []byte{}
	}
	buf := make([]byte, 2+len(text))
	binary.BigEndian.PutUint16(buf, uint16(closeCode))
	copy(buf[2:], text)
	return buf
}

// Write methods

func (s *Storm) writeFatal(err error) error {
	s.writeErrMu.Lock()
	if s.writeErr == nil {
		s.writeErr = err
	}
	s.writeErrMu.Unlock()
	return err
}

func (s *Storm) write(frameType int, deadline time.Time, buf0, buf1 []byte) error {
	<-s.mu
	defer func() { s.mu <- struct{}{} }()

	s.writeErrMu.Lock()
	err := s.writeErr
	s.writeErrMu.Unlock()
	if err != nil {
		return err
	}

	if err := s.RawConn.SetWriteDeadline(deadline); err != nil {
		return s.writeFatal(err)
	}
	if len(buf1) == 0 {
		_, err = s.RawConn.Write(buf0)
	} else {
		err = s.writeBufs(buf0, buf1)
	}
	if err != nil {
		return s.writeFatal(err)
	}
	if frameType == CloseMessage {
		_ = s.writeFatal(ErrCloseSent)
	}
	return nil
}

func (s *Storm) writeBufs(bufs ...[]byte) error {
	b := net.Buffers(bufs)
	_, err := b.WriteTo(s.RawConn)
	return err
}
