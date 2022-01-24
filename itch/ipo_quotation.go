/*
 * Copyright (c) 2022 Mark Edward Winter
 */
package itch

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

type ReleaseQualifier uint8

const (
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
	Price          uint32
}

func MakeIpoQuotation(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return IpoQuotation{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Stock:          strings.TrimSpace(string(data[11:19])),
		Qualifier:      ReleaseQualifier(data[23]),
		Price:          binary.BigEndian.Uint32(data[24:28]),
	}
}

func (i IpoQuotation) String() string {
	return fmt.Sprintf("[IPO Quotation]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Stock: %v\n"+
		"Qualifier: %v\n"+
		"Price: %v\n",
		i.StockLocate, i.TrackingNumber, i.Timestamp,
		i.Stock, i.Qualifier, float64(i.Price)/10000,
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
