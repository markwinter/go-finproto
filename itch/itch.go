/*
 * Copyright (c) 2022 Mark Edward Winter
 */

// Package itch implements the Nasdaq ITCH 5.0 protocol
package itch

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

const (
	MESSAGE_SYSTEM_EVENT         uint8 = 'S'
	MESSAGE_STOCK_DIRECTORY      uint8 = 'R'
	MESSAGE_STOCK_TRADING_ACTION uint8 = 'H'
	MESSAGE_REG_SHO              uint8 = 'Y'
	MESSAGE_PARTICIPANT_POSITION uint8 = 'L'
	MESSAGE_MCWB_LEVEL           uint8 = 'V'
	MESSAGE_MCWB_STATUS          uint8 = 'W'
	MESSAGE_IPO_QUOTATION        uint8 = 'K'
	MESSAGE_LULD_COLLAR          uint8 = 'J'
	MESSAGE_OPERATIONAL_HALT     uint8 = 'h'
	MESSAGE_ORDER_ADD            uint8 = 'A'
	MESSAGE_ORDER_ADD_ATTRIBUTED uint8 = 'F'
	MESSAGE_ORDER_EXECUTED       uint8 = 'E'
	MESSAGE_ORDER_EXECUTED_PRICE uint8 = 'C'
	MESSAGE_ORDER_CANCEL         uint8 = 'X'
	MESSAGE_ORDER_DELETE         uint8 = 'D'
	MESSAGE_ORDER_REPLACE        uint8 = 'U'
	MESSAGE_TRADE_NON_CROSS      uint8 = 'P'
	MESSAGE_TRADE_CROSS          uint8 = 'Q'
	MESSAGE_TRADE_BROKEN         uint8 = 'B'
	MESSAGE_NOII                 uint8 = 'I'
	MESSAGE_RPII                 uint8 = 'N'
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
	case MESSAGE_SYSTEM_EVENT:
		return MakeSystemEvent(data)
	case MESSAGE_STOCK_DIRECTORY:
		return MakeStockDirectory(data)
	case MESSAGE_STOCK_TRADING_ACTION:
		return MakeStockTradingAction(data)
	case MESSAGE_REG_SHO:
		return MakeRegSho(data)
	case MESSAGE_PARTICIPANT_POSITION:
		return MakeParticipantPosition(data)
	case MESSAGE_MCWB_LEVEL:
		return MakeMcwbLevel(data)
	case MESSAGE_MCWB_STATUS:
		return MakeMcwbStatus(data)
	case MESSAGE_IPO_QUOTATION:
		return MakeIpoQuotation(data)
	case MESSAGE_LULD_COLLAR:
		return MakeLuldCollar(data)
	case MESSAGE_OPERATIONAL_HALT:
		return MakeOperationalHalt(data)
	case MESSAGE_ORDER_ADD:
		return MakeOrderAdd(data)
	case MESSAGE_ORDER_ADD_ATTRIBUTED:
		return MakeOrderAddAttributed(data)
	case MESSAGE_ORDER_EXECUTED:
		return MakeOrderExecuted(data)
	case MESSAGE_ORDER_EXECUTED_PRICE:
		return MakeOrderExecutedPrice(data)
	case MESSAGE_ORDER_CANCEL:
		return MakeOrderCancel(data)
	case MESSAGE_ORDER_DELETE:
		return MakeOrderDelete(data)
	case MESSAGE_ORDER_REPLACE:
		return MakeOrderReplace(data)
	case MESSAGE_TRADE_NON_CROSS:
		return MakeTradeNonCross(data)
	case MESSAGE_TRADE_CROSS:
		return MakeTradeCross(data)
	case MESSAGE_TRADE_BROKEN:
		return MakeTradeBroken(data)
	case MESSAGE_NOII:
		return MakeNoii(data)
	case MESSAGE_RPII:
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
