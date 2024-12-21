package itch

import (
	"bytes"
	"encoding/binary"
	"unicode"

	"github.com/quagmt/udecimal"
)

func bytesToPrice(data []byte, precision uint8) (udecimal.Decimal, error) {
	// Prices are integer fields, supplied with an associated precision.
	// When converted to a decimal format, prices are in fixed point format,
	// where the precision defines the number of decimal places.
	// For example, a field flagged as Price (4) has an implied 4 decimal places.
	// The maximum value of price (4) in TotalViewITCH is 200,000.0000 (decimal, 77359400 hex)
	data = bytes.TrimLeftFunc(data, unicode.IsSpace)
	value := binary.BigEndian.Uint64(append(make([]byte, 8-len(data)), data...))
	return udecimal.NewFromUint64(value, precision)
}

func priceToBytes(price udecimal.Decimal, precision uint8) ([]byte, error) {
	price, _ = price.Div(udecimal.MustFromInt64(1, precision))
	p, _ := price.Int64()

	bytes := make([]byte, precision)
	for i := 0; i < int(precision); i++ {
		shift := uint((int(precision) - 1 - i) * 8)
		bytes[i] = byte(p >> shift)
	}

	return bytes, nil
}
