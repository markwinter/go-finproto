/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"encoding/binary"
	"fmt"
	"time"
)

type EventCode uint8

const (
	EVENT_START_MESSAGES EventCode = 'O'
	EVENT_START_HOURS    EventCode = 'S'
	EVENT_START_MARKET   EventCode = 'Q'
	EVENT_END_MARKET     EventCode = 'M'
	EVENT_END_HOURS      EventCode = 'E'
	EVENT_END_MESSAGES   EventCode = 'C'
)

type SystemEvent struct {
	Timestamp      time.Duration
	StockLocate    uint16
	TrackingNumber uint16
	EventCode      EventCode
}

func (e SystemEvent) Type() uint8 {
	return MESSAGE_SYSTEM_EVENT
}

func (e SystemEvent) Bytes() []byte {
	data := make([]byte, systemEventSize)

	data[0] = MESSAGE_SYSTEM_EVENT
	binary.BigEndian.PutUint16(data[1:3], 0)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(e.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], e.TrackingNumber)

	data[11] = byte(e.EventCode)

	return data
}

func ParseSystemEvent(data []byte) SystemEvent {
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

func MakeSystemEvent(timestamp time.Duration, trackingNumber uint16, eventCode EventCode) SystemEvent {
	return SystemEvent{
		Timestamp:      timestamp,
		StockLocate:    0, // StockLocate for System Event is always 0
		TrackingNumber: trackingNumber,
		EventCode:      eventCode,
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
