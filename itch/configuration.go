package itch

// Configuration contains settings for adjusting how messages are parsed
type Configuration struct {
	// Set which message types to parse
	MessageTypes []byte
	// Maximum amount of messages to parse
	MaxMessages int
}
