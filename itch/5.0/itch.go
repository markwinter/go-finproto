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
	"slices"
	"time"
)

const (
	MESSAGE_STOCK_DIRECTORY      uint8 = 'R'
	MESSAGE_STOCK_TRADING_ACTION uint8 = 'H'
	MESSAGE_REG_SHO              uint8 = 'Y'
	MESSAGE_PARTICIPANT_POSITION uint8 = 'L'
	MESSAGE_MWCB_LEVEL           uint8 = 'V'
	MESSAGE_MWCB_STATUS          uint8 = 'W'
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

	// Message lengths are fixed sized.
	stockDirectorySize      = 39
	stockTradingActionSize  = 25
	regShoSize              = 20
	participantPositionSize = 26
	mwcbLevelSize           = 35
	mwcbStatusSize          = 12
	luldSize                = 35
	operationalHaltSize     = 21
	orderAddSize            = 36
	orderAddAttrSize        = 40
	orderExecutedSize       = 31
	orderExecutedPriceSize  = 36
	orderCancelSize         = 23
	orderDeleteSize         = 19
	orderReplaceSize        = 35
	tradeNonCrossSize       = 44
	tradeCrossSize          = 40
	tradeBrokenSize         = 19
	noiiSize                = 50
	rpiiSize                = 20
)

type ItchMessage interface {
	Bytes() []byte
	Type() uint8
}

// ParseFile parses ITCH messages from an uncompressed file. It uses ParseReader internally
func ParseFile(path string, config Configuration) ([]ItchMessage, error) {
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

	log.Printf("Using buffer size: %v\n", reader.Size())

	return ParseReader(reader, config)
}

// ParseReader parses ITCH messages from a bufio.Reader
func ParseReader(reader *bufio.Reader, config Configuration) ([]ItchMessage, error) {
	messages := []ItchMessage{}

	start := time.Now()

	for {
		if config.MaxMessages > 0 && len(messages) >= config.MaxMessages {
			break
		}

		var msgLength int

		if config.LengthFieldPrefixed {
			msgLengthBuffer, err := reader.Peek(2)
			if err == io.EOF {
				break
			}
			if err != nil {
				return messages, err
			}

			_, err = reader.Discard(2)
			if err != nil {
				return messages, err
			}

			msgLength = int(uint16(msgLengthBuffer[1]) | uint16(msgLengthBuffer[0])<<8)

		} else {
			msgTypeBuffer, err := reader.Peek(1)
			if err == io.EOF {
				break
			}
			if err != nil {
				return messages, err
			}

			msgLength = getMessageSize(msgTypeBuffer[0])
		}

		data := make([]byte, msgLength)

		n, err := io.ReadFull(reader, data)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return messages, err
		}

		// If user configured MessageTypes then only parse messages they want
		if len(config.MessageTypes) != 0 {
			if !slices.Contains(config.MessageTypes, data[0]) {
				continue
			}
		}

		messages = append(messages, parseData(data[0], data))
	}

	elapsed := time.Since(start)
	log.Printf("Parsed %d messages in %s", len(messages), elapsed)
	log.Printf("Parse rate: %.2f messages/s", float64(len(messages))/elapsed.Seconds())

	return messages, nil
}

// ParseMany parses multiple ITCH messages from byte data already loaded into memory
func ParseMany(data []byte, config Configuration) ([]ItchMessage, error) {
	messages := []ItchMessage{}

	start := time.Now()

	dp := 0

	for {
		if dp >= len(data) {
			// Reached end of data
			break
		}

		if config.MaxMessages > 0 && len(messages) >= config.MaxMessages {
			break
		}

		var msgLength int
		if config.LengthFieldPrefixed {
			msgLength = int(uint16(data[dp+1]) | uint16(data[dp])<<8)
		} else {
			msgLength = getMessageSize(data[dp])
		}

		startOfMessage := dp
		endOfMessage := startOfMessage + msgLength

		dp += int(msgLength)

		// If user configured MessageTypes then only parse messages they want
		if len(config.MessageTypes) != 0 {
			if !slices.Contains(config.MessageTypes, data[startOfMessage]) {
				continue
			}
		}

		messages = append(messages, parseData(data[startOfMessage], data[startOfMessage:endOfMessage]))
	}

	elapsed := time.Since(start)

	log.Printf("Parsed %d messages in %s", len(messages), elapsed)
	log.Printf("Parse rate: %.2f messages/s", float64(len(messages))/elapsed.Seconds())

	return messages, nil
}

// Parse will parse a single ITCH message - it should not have a length field prefixed, just give the actual ITCH message
func Parse(data []byte) ItchMessage {
	return parseData(data[0], data)
}

func getMessageSize(msgType byte) int {
	switch msgType {
	case MESSAGE_SYSTEM_EVENT:
		return systemEventSize
	case MESSAGE_STOCK_DIRECTORY:
		return stockDirectorySize
	case MESSAGE_STOCK_TRADING_ACTION:
		return stockTradingActionSize
	case MESSAGE_REG_SHO:
		return regShoSize
	case MESSAGE_PARTICIPANT_POSITION:
		return participantPositionSize
	case MESSAGE_MWCB_LEVEL:
		return mwcbLevelSize
	case MESSAGE_MWCB_STATUS:
		return mwcbStatusSize
	case MESSAGE_IPO_QUOTATION:
		return ipoQuotationSize
	case MESSAGE_LULD_COLLAR:
		return luldSize
	case MESSAGE_OPERATIONAL_HALT:
		return operationalHaltSize
	case MESSAGE_ORDER_ADD:
		return orderAddSize
	case MESSAGE_ORDER_ADD_ATTRIBUTED:
		return orderAddAttrSize
	case MESSAGE_ORDER_EXECUTED:
		return orderExecutedSize
	case MESSAGE_ORDER_EXECUTED_PRICE:
		return orderExecutedPriceSize
	case MESSAGE_ORDER_CANCEL:
		return orderCancelSize
	case MESSAGE_ORDER_DELETE:
		return orderDeleteSize
	case MESSAGE_ORDER_REPLACE:
		return orderReplaceSize
	case MESSAGE_TRADE_NON_CROSS:
		return tradeNonCrossSize
	case MESSAGE_TRADE_CROSS:
		return tradeCrossSize
	case MESSAGE_TRADE_BROKEN:
		return tradeBrokenSize
	case MESSAGE_NOII:
		return noiiSize
	case MESSAGE_RPII:
		return rpiiSize
	default:
		return 0
	}
}

func parseData(msgType byte, data []byte) ItchMessage {
	switch msgType {
	case MESSAGE_SYSTEM_EVENT:
		return ParseSystemEvent(data)
	case MESSAGE_STOCK_DIRECTORY:
		return ParseStockDirectory(data)
	case MESSAGE_STOCK_TRADING_ACTION:
		return ParseStockTradingAction(data)
	case MESSAGE_REG_SHO:
		return ParseRegSho(data)
	case MESSAGE_PARTICIPANT_POSITION:
		return ParseParticipantPosition(data)
	case MESSAGE_MWCB_LEVEL:
		return ParseMwcbLevel(data)
	case MESSAGE_MWCB_STATUS:
		return ParseMwcbStatus(data)
	case MESSAGE_IPO_QUOTATION:
		return ParseIpoQuotation(data)
	case MESSAGE_LULD_COLLAR:
		return ParseLuldCollar(data)
	case MESSAGE_OPERATIONAL_HALT:
		return ParseOperationalHalt(data)
	case MESSAGE_ORDER_ADD:
		return ParseOrderAdd(data)
	case MESSAGE_ORDER_ADD_ATTRIBUTED:
		return ParseOrderAddAttributed(data)
	case MESSAGE_ORDER_EXECUTED:
		return ParseOrderExecuted(data)
	case MESSAGE_ORDER_EXECUTED_PRICE:
		return ParseOrderExecutedPrice(data)
	case MESSAGE_ORDER_CANCEL:
		return ParseOrderCancel(data)
	case MESSAGE_ORDER_DELETE:
		return ParseOrderDelete(data)
	case MESSAGE_ORDER_REPLACE:
		return ParseOrderReplace(data)
	case MESSAGE_TRADE_NON_CROSS:
		return ParseTradeNonCross(data)
	case MESSAGE_TRADE_CROSS:
		return ParseTradeCross(data)
	case MESSAGE_TRADE_BROKEN:
		return ParseTradeBroken(data)
	case MESSAGE_NOII:
		return ParseNoii(data)
	case MESSAGE_RPII:
		return ParseRpii(data)
	default:
		return nil
	}
}
