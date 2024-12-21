package itch

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMakeAndParseIpoQuotation(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/London")
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	timeSinceMidnight := now.Sub(midnight).Truncate(time.Second)

	manual := IpoQuotation{
		StockLocate:    0,
		TrackingNumber: 1,
		Timestamp:      timeSinceMidnight,
		Stock:          "AAPL    ",
		ReleaseTime:    timeSinceMidnight,
		Qualifier:      QUALIFIER_ANTICIPATED,
		Price:          21,
	}

	m := MakeIpoQuotation(0, 1, timeSinceMidnight, "AAPL", timeSinceMidnight, QUALIFIER_ANTICIPATED, 21)

	if !cmp.Equal(manual, m) {
		t.Errorf("created event and manual struct not equal:\n%v", cmp.Diff(manual, m))
	}

	parsedEvent, err := ParseIpoQuotation(m.Bytes())
	if err != nil {
		t.Errorf("error parsing event: %v", err)
	}

	// Ignoring Stock field because the byte representation contains right-padded spacing, whilst the Struct has it trimmed
	if !cmp.Equal(m, parsedEvent, cmpopts.IgnoreFields(IpoQuotation{}, "Stock")) {
		t.Errorf("parsed event and original event are not equal:\n%v", cmp.Diff(m, parsedEvent))
	}
}
