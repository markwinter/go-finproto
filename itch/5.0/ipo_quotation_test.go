package itch

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/quagmt/udecimal"
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
		Stock:          "AAPL",
		ReleaseTime:    timeSinceMidnight,
		Qualifier:      QUALIFIER_ANTICIPATED,
		Price:          udecimal.MustFromFloat64(21.45),
	}

	m := MakeIpoQuotation(0, 1, timeSinceMidnight, "AAPL", timeSinceMidnight, QUALIFIER_ANTICIPATED, udecimal.MustParse("21.45"))

	if !cmp.Equal(manual, m) {
		t.Errorf("created event and manual struct not equal:\n%v", cmp.Diff(manual, m))
	}

	parsedEvent, err := ParseIpoQuotation(m.Bytes())
	if err != nil {
		t.Errorf("error parsing event: %v", err)
	}

	if !cmp.Equal(m, parsedEvent) {
		t.Errorf("parsed event and original event are not equal:\n%v", cmp.Diff(m, parsedEvent))
	}
}

func TestParseIpoQuotation(t *testing.T) {
	bdtxTimestamp, _ := time.ParseDuration("8h25m39.142161035s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    IpoQuotation
		wantErr bool
	}{
		{
			name: "BDTX IPO",
			args: args{
				// Real data taken from an ITCH file
				data: []byte{75, 0, 0, 0, 0, 27, 151, 225, 202, 146, 139, 66, 68, 84, 88, 32, 32, 32, 32, 0, 0, 142, 248, 65, 0, 2, 230, 48},
			},
			want: IpoQuotation{
				StockLocate:    0,
				TrackingNumber: 0,
				Stock:          "BDTX",
				Timestamp:      bdtxTimestamp,
				ReleaseTime:    36600 * time.Second,
				Qualifier:      QUALIFIER_ANTICIPATED,
				Price:          udecimal.MustParse("19.0000"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIpoQuotation(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIpoQuotation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}
