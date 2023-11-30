package itch

import (
	"testing"
	"time"
)

func TestMakeAndParse(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/London")
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	timeSinceMidnight := now.Sub(midnight)

	manual := SystemEvent{
		StockLocate:    0,
		Timestamp:      timeSinceMidnight,
		TrackingNumber: 1,
		EventCode:      EVENT_START_HOURS,
	}

	s := MakeSystemEvent(timeSinceMidnight, 1, EVENT_START_HOURS)

	if manual.EventCode != s.EventCode {
		t.Errorf("EventCode incorrect, got: %s, want %s", s.EventCode, manual.EventCode)
	}
	if manual.StockLocate != s.StockLocate {
		t.Errorf("EventCode incorrect, got: %d, want: %d", s.StockLocate, manual.StockLocate)
	}
	if manual.Timestamp != s.Timestamp {
		t.Errorf("EventCode incorrect, got: %s, want: %s", s.Timestamp, manual.Timestamp)
	}
	if manual.TrackingNumber != s.TrackingNumber {
		t.Errorf("EventCode incorrect, got: %d, want: %d", s.TrackingNumber, manual.TrackingNumber)
	}

	parsedEvent := ParseSystemEvent(s.Bytes())

	if parsedEvent.EventCode != s.EventCode {
		t.Errorf("EventCode incorrect, got: %s, want %s", parsedEvent.EventCode, s.EventCode)
	}
	if parsedEvent.StockLocate != s.StockLocate {
		t.Errorf("EventCode incorrect, got: %d, want: %d", parsedEvent.StockLocate, s.StockLocate)
	}
	if parsedEvent.Timestamp != s.Timestamp {
		t.Errorf("EventCode incorrect, got: %s, want: %s", parsedEvent.Timestamp, s.Timestamp)
	}
	if parsedEvent.TrackingNumber != s.TrackingNumber {
		t.Errorf("EventCode incorrect, got: %d, want: %d", parsedEvent.TrackingNumber, s.TrackingNumber)
	}
}
