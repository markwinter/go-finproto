/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/quagmt/udecimal"
)

type ReleaseQualifier uint8

const (
	MESSAGE_IPO_QUOTATION uint8 = 'K'

	ipoQuotationSize = 28

	QUALIFIER_ANTICIPATED       ReleaseQualifier = 'A'
	QUALIFER_CANCELED_POSTPONED ReleaseQualifier = 'C'
)

type IpoQuotation struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	Stock          string
	ReleaseTime    time.Duration
	Qualifier      ReleaseQualifier
	Price          udecimal.Decimal // Price (4)
}

func (i IpoQuotation) Type() uint8 {
	return MESSAGE_IPO_QUOTATION
}

func (i IpoQuotation) Bytes() []byte {
	data := make([]byte, ipoQuotationSize)

	data[0] = MESSAGE_IPO_QUOTATION
	binary.BigEndian.PutUint16(data[1:3], i.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(i.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], i.TrackingNumber)

	copy(data[11:19], []byte(fmt.Sprintf("%-8s", i.Stock)))

	binary.BigEndian.PutUint32(data[19:23], uint32(i.ReleaseTime.Seconds()))

	data[23] = byte(i.Qualifier)

	price, _ := priceToBytes(i.Price, 4)
	copy(data[24:28], price)

	return data
}

func MakeIpoQuotation(stockLocate, trackingNumber uint16, timestamp time.Duration, stock string, releaseTime time.Duration, qualifier ReleaseQualifier, price udecimal.Decimal) IpoQuotation {
	return IpoQuotation{
		StockLocate:    stockLocate,
		TrackingNumber: trackingNumber,
		Timestamp:      timestamp,
		Stock:          stock,
		ReleaseTime:    releaseTime,
		Qualifier:      qualifier,
		Price:          price,
	}
}

func ParseIpoQuotation(data []byte) (IpoQuotation, error) {
	if len(data) != ipoQuotationSize {
		return IpoQuotation{}, NewInvalidPacketSize(ipoQuotationSize, len(data))
	}

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	stock := strings.TrimSpace(string(data[11:19]))

	releaseTime := binary.BigEndian.Uint32(data[19:23])

	price, _ := bytesToPrice(data[24:28], 4)

	return IpoQuotation{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Stock:          stock,
		ReleaseTime:    time.Duration(uint64(releaseTime) * uint64(time.Second)),
		Qualifier:      ReleaseQualifier(data[23]),
		Price:          price,
	}, nil
}

func (i IpoQuotation) String() string {
	return fmt.Sprintf("[IPO Quotation]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Stock: %v\n"+
		"Release Time: %ds\n"+
		"Qualifier: %v\n"+
		"Price: %v\n",
		i.StockLocate, i.TrackingNumber, i.Timestamp,
		i.Stock, int64(i.ReleaseTime.Seconds()), i.Qualifier, i.Price.String(),
	)
}

func (r ReleaseQualifier) String() string {
	switch r {
	case QUALIFIER_ANTICIPATED:
		return "Anticipated"
	case QUALIFER_CANCELED_POSTPONED:
		return "IPO release cancelled/postponed"
	}

	return "Unknown ReleaseQualifier"
}
