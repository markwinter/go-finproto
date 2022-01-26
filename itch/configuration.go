/*
 * Copyright (c) 2022 Mark Edward Winter
 */
package itch

const (
	OneGB = 1 << (10 * 3)
)

// Configuration contains settings for adjusting how messages are parsed
type Configuration struct {
	// Set which message types to parse
	MessageTypes []byte
	// Maximum amount of messages to parse
	MaxMessages int
	// Set buffer size for io.reader when using ParseFile
	ReadBufferSize uint64
}
