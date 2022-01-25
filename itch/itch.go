/*
 * Copyright (c) 2022 Mark Edward Winter
 */
package itch

import (
	"encoding/binary"
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

	return ParseReader(file, config)
}

// ParseReader parses ITCH messages from an io.Reader
func ParseReader(reader io.Reader, config Configuration) ([]Message, error) {
	messages := []Message{}
	messageCount := 0

	var msgLength uint16

	start := time.Now()

	for {
		if config.MaxMessages > 0 && messageCount >= config.MaxMessages {
			break
		}

		err := binary.Read(reader, binary.BigEndian, &msgLength)
		if err != nil {
			// Reached the end of data
			if err == io.ErrUnexpectedEOF || err == io.EOF {
				break
			}
			return messages, err
		}

		buffer := make([]byte, msgLength)
		_, err = reader.Read(buffer)
		if err != nil {
			return messages, err
		}

		messageCount++

		// If user configured MessageTypes then only parse messages they want
		if len(config.MessageTypes) != 0 {
			if !contains(config.MessageTypes, buffer[0]) {
				continue
			}
		}

		messages = append(messages, makeMessage(buffer[0], buffer))
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

	msgLength := uint16(0)

	start := time.Now()

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
		return MakeAddOrder(data)
	case 'F':
		return MakeAddOrderAttributed(data)
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
