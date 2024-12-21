package itch

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestMakeAndParseSystemEvent(t *testing.T) {
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

	if !cmp.Equal(manual, s) {
		t.Errorf("created event and manual struct not equal:\n%v\n%v", manual, s)
	}

	parsedEvent, err := ParseSystemEvent(s.Bytes())
	if err != nil {
		t.Errorf("got error from ParseSystemEvent: %s", err)
	}

	if !cmp.Equal(parsedEvent, s) {
		t.Errorf("parsedEvent and created event not equal:\n%v\n%v", s, parsedEvent)
	}
}
