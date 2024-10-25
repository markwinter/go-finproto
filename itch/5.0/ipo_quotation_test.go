package itch

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestMakeAndParseIpoQuotation(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/London")
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	timeSinceMidnight := now.Sub(midnight)

	manual := IpoQuotation{
		StockLocate:    0,
		TrackingNumber: 1,
		Timestamp:      timeSinceMidnight,
		Stock:          "AAPL    ",
		ReleaseTime:    timeSinceMidnight,
		Qualifier:      QUALIFIER_ANTICIPATED,
		Price:          21.00,
	}

	m := MakeIpoQuotation(0, 1, timeSinceMidnight, "AAPL", timeSinceMidnight, QUALIFIER_ANTICIPATED, 20)

	if !cmp.Equal(manual, m) {
		t.Errorf("created event and manual struct not equal:\n%v\n%v", manual, m)
	}

	parsedEvent := ParseIpoQuotation(m.Bytes())

	if !cmp.Equal(m, parsedEvent) {
		t.Errorf("parsed event and original event are not equal:\n%v\n%v", m, parsedEvent)
	}
}
