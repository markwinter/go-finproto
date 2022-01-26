/*
 * Copyright (c) 2022 Mark Edward Winter
 */
package itch

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

type Message interface{}

// ParseFile parses ITCH messages from an uncompressed file. It uses ParseReader internally
func ParseFile(path string, config Configuration) ([]Message, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var reader *bufio.Reader
	if config.ReadBufferSize > 0 {
		reader = bufio.NewReaderSize(file, int(config.ReadBufferSize))
	} else {
		reader = bufio.NewReader(file)
	}

	log.Printf("Using buffer size: %v", reader.Size())
	log.Printf("Buffer size: %v", reader.Buffered())

	return ParseReader(reader, config)
}

// ParseReader parses ITCH messages from a bufio.Reader
func ParseReader(reader *bufio.Reader, config Configuration) ([]Message, error) {
	messages := []Message{}
	messageCount := 0

	start := time.Now()

	for {
		if config.MaxMessages > 0 && messageCount >= config.MaxMessages {
			break
		}

		msgLengthBuffer, err := reader.Peek(2)
		if err == io.EOF {
			break
		}
		if err != nil {
			return messages, err
		}
		reader.Discard(2)

		msgLength := uint16(msgLengthBuffer[1]) | uint16(msgLengthBuffer[0])<<8

		data, err := reader.Peek(int(msgLength))
		if err == io.EOF {
			break
		}
		if err != nil {
			return messages, err
		}
		reader.Discard(int(msgLength))

		messageCount++

		// If user configured MessageTypes then only parse messages they want
		if len(config.MessageTypes) != 0 {
			if !contains(config.MessageTypes, data[0]) {
				continue
			}
		}

		messages = append(messages, makeMessage(data[0], data))
	}

	elapsed := time.Since(start)
	log.Printf("Parsed %d messages in %s", messageCount, elapsed)
	log.Printf("Parse rate: %.2f messages/s", float64(messageCount)/elapsed.Seconds())

	return messages, nil
}

// ParseMany parses multiple ITCH messages from byte data already loaded into memory
func ParseMany(data []byte, config Configuration) ([]Message, error) {
	messages := []Message{}
	messageCount := 0

	start := time.Now()

	msgLength := uint16(0)
	dp := 0

	for {
		if config.MaxMessages > 0 && messageCount >= config.MaxMessages {
			break
		}

		dp += int(msgLength)
		if dp >= len(data) {
			// Reached end of data
			break
		}

		// This is quicker than using binary.BigEndian.Uint16 in this loop
		msgLength = uint16(data[dp+1]) | uint16(data[dp])<<8
		dp += 2

		messageCount++

		// If user configured MessageTypes then only parse messages they want
		if len(config.MessageTypes) != 0 {
			if !contains(config.MessageTypes, data[dp]) {
				continue
			}
		}

		messages = append(messages, makeMessage(data[dp], data[dp:dp+int(msgLength)]))
	}

	elapsed := time.Since(start)
	log.Printf("Parsed %d messages in %s", messageCount, elapsed)
	log.Printf("Parse rate: %.2f messages/s", float64(messageCount)/elapsed.Seconds())

	return messages, nil
}

// Parse will parse a single ITCH message
func Parse(data []byte) Message {
	return makeMessage(data[2], data)
}

func makeMessage(msgType byte, data []byte) Message {
	switch msgType {
	case 'S':
		return MakeSystemEvent(data)
	case 'R':
		return MakeStockDirectory(data)
	case 'H':
		return MakeStockTradingAction(data)
	case 'Y':
		return MakeRegSho(data)
	case 'L':
		return MakeParticipantPosition(data)
	case 'V':
		return MakeMcwbLevel(data)
	case 'W':
		return MakeMcwbStatus(data)
	case 'K':
		return MakeIpoQuotation(data)
	case 'J':
		return MakeLuldCollar(data)
	case 'h':
		return MakeOperationalHalt(data)
	case 'A':
		return MakeOrderAdd(data)
	case 'F':
		return MakeOrderAddAttributed(data)
	case 'E':
		return MakeOrderExecuted(data)
	case 'C':
		return MakeOrderExecutedPrice(data)
	case 'X':
		return MakeOrderCancel(data)
	case 'D':
		return MakeOrderDelete(data)
	case 'U':
		return MakeOrderReplace(data)
	case 'P':
		return MakeTradeNonCross(data)
	case 'Q':
		return MakeTradeCross(data)
	case 'B':
		return MakeTradeBroken(data)
	case 'I':
		return MakeNoii(data)
	case 'N':
		return MakeRpii(data)
	default:
		return nil
	}
}

func contains(l []byte, x byte) bool {
	for i := 0; i < len(l); i++ {
		if l[i] == x {
			return true
		}
	}

	return false
}
