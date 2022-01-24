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

		switch buffer[0] {
		case 'S':
			messages = append(messages, MakeSystemEvent(buffer))
		case 'R':
			messages = append(messages, MakeStockDirectory(buffer))
		case 'H':
			messages = append(messages, MakeStockTradingAction(buffer))
		case 'Y':
			messages = append(messages, MakeRegSho(buffer))
		case 'L':
			messages = append(messages, MakeParticipantPosition(buffer))
		case 'V':
			messages = append(messages, MakeMcwbLevel(buffer))
		case 'W':
			messages = append(messages, MakeMcwbStatus(buffer))
		case 'K':
			messages = append(messages, MakeIpoQuotation(buffer))
		case 'J':
			messages = append(messages, MakeLuldCollar(buffer))
		case 'h':
			messages = append(messages, MakeOperationalHalt(buffer))
		case 'A':
			messages = append(messages, MakeAddOrder(buffer))
		case 'F':
			messages = append(messages, MakeAddOrderAttributed(buffer))
		default:
			continue
		}
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

		switch data[dp] {
		case 'S':
			messages = append(messages, MakeSystemEvent(data[dp:dp+int(msgLength)]))
		case 'R':
			messages = append(messages, MakeStockDirectory(data[dp:dp+int(msgLength)]))
		case 'H':
			messages = append(messages, MakeStockTradingAction(data[dp:dp+int(msgLength)]))
		case 'Y':
			messages = append(messages, MakeRegSho(data[dp:dp+int(msgLength)]))
		case 'L':
			messages = append(messages, MakeParticipantPosition(data[dp:dp+int(msgLength)]))
		case 'V':
			messages = append(messages, MakeMcwbLevel(data[dp:dp+int(msgLength)]))
		case 'W':
			messages = append(messages, MakeMcwbStatus(data[dp:dp+int(msgLength)]))
		case 'K':
			messages = append(messages, MakeIpoQuotation(data[dp:dp+int(msgLength)]))
		case 'J':
			messages = append(messages, MakeLuldCollar(data[dp:dp+int(msgLength)]))
		case 'h':
			messages = append(messages, MakeOperationalHalt(data[dp:dp+int(msgLength)]))
		case 'A':
			messages = append(messages, MakeAddOrder(data[dp:dp+int(msgLength)]))
		case 'F':
			messages = append(messages, MakeAddOrderAttributed(data[dp:dp+int(msgLength)]))
		default:
			continue
		}
	}

	elapsed := time.Since(start)
	log.Printf("Parsed %d messages in %s", messageCount, elapsed)
	log.Printf("Parse rate: %.2f messages/s", float64(messageCount)/elapsed.Seconds())

	return messages, nil
}

func contains(l []byte, x byte) bool {
	for i := 0; i < len(l); i++ {
		if l[i] == x {
			return true
		}
	}

	return false
}
