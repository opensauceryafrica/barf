package websocket

import "time"

const (
	maxFrameHeaderSize         = 2 + 8 + 4 // Fixed header + length + mask
	maxControlFramePayloadSize = 125

	// default buffer sizes
	defaultReadBufferSize  = 4096
	defaultWriteBufferSize = 4096

	// Frame header byte 0 bits from Section 5.2 of RFC 6455
	finalBit = 1 << 7
	rsv1Bit  = 1 << 6
	rsv2Bit  = 1 << 5
	rsv3Bit  = 1 << 4

	// Frame header byte 1 bits from Section 5.2 of RFC 6455
	maskBit = 1 << 7

	writeWait = time.Second

	continuationFrame = 0
	noFrame           = -1
)
