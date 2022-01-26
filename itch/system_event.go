/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"encoding/binary"
	"fmt"
	"time"
)

type SystemEvent struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	EventCode      EventCode
}

type EventCode uint8

const (
	EVENT_START_MESSAGES EventCode = 'O'
	EVENT_START_HOURS    EventCode = 'S'
	EVENT_START_MARKET   EventCode = 'Q'
	EVENT_END_MARKET     EventCode = 'M'
	EVENT_END_HOURS      EventCode = 'E'
	EVENT_END_MESSAGES   EventCode = 'C'
)

func MakeSystemEvent(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])
	event := EventCode(data[11])

	return SystemEvent{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		EventCode:      event,
	}
}

func (e SystemEvent) String() string {
	return fmt.Sprintf("[System Event]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Event Code: %v\n",
		e.StockLocate, e.TrackingNumber, e.Timestamp, e.EventCode)
}

func (e EventCode) String() string {
	switch e {
	case EVENT_START_MESSAGES:
		return "Start of Messages"
	case EVENT_START_HOURS:
		return "Start of Hours"
	case EVENT_START_MARKET:
		return "Start of Market"
	case EVENT_END_MARKET:
		return "End of Market"
	case EVENT_END_HOURS:
		return "End of Hours"
	case EVENT_END_MESSAGES:
		return "End of Messages"
	}

	return "Unknown EventCode"
}
