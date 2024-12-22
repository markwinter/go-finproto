package itch

import (
	"reflect"
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

func TestParseSystemEvent(t *testing.T) {
	// These test contain real data taken from ITCH file

	startMessagesTimestamp, _ := time.ParseDuration("3h2m33.404452051s")
	startHoursTimestamp, _ := time.ParseDuration("4h0m0.000245942s")
	startMarketTimestamp, _ := time.ParseDuration("9h30m0.000073183s")
	endMarketTimestamp, _ := time.ParseDuration("16h0m0.00007852s")
	endHoursTimestamp, _ := time.ParseDuration("20h0m0.000021649s")
	endMessagesTimestamp, _ := time.ParseDuration("20h5m0.000039817s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    SystemEvent
		wantErr bool
	}{
		{
			name: "Start of Messages",
			args: args{
				data: []byte{83, 0, 0, 0, 0, 9, 246, 73, 200, 12, 211, 79},
			},
			want: SystemEvent{
				StockLocate:    0,
				TrackingNumber: 0,
				EventCode:      EVENT_START_MESSAGES,
				Timestamp:      startMessagesTimestamp,
			},
			wantErr: false,
		},
		{
			name: "Start of Hours",
			args: args{
				data: []byte{83, 0, 0, 0, 0, 13, 24, 194, 230, 64, 182, 83},
			},
			want: SystemEvent{
				StockLocate:    0,
				TrackingNumber: 0,
				EventCode:      EVENT_START_HOURS,
				Timestamp:      startHoursTimestamp,
			},
			wantErr: false,
		},
		{
			name: "Start of Market",
			args: args{
				data: []byte{83, 0, 0, 0, 0, 31, 26, 206, 219, 13, 223, 81},
			},
			want: SystemEvent{
				StockLocate:    0,
				TrackingNumber: 0,
				EventCode:      EVENT_START_MARKET,
				Timestamp:      startMarketTimestamp,
			},
			wantErr: false,
		},
		{
			name: "End of Market",
			args: args{
				data: []byte{83, 0, 0, 0, 0, 52, 99, 11, 139, 50, 184, 77},
			},
			want: SystemEvent{
				StockLocate:    0,
				TrackingNumber: 0,
				EventCode:      EVENT_END_MARKET,
				Timestamp:      endMarketTimestamp,
			},
			wantErr: false,
		},
		{
			name: "End of Hours",
			args: args{
				data: []byte{83, 0, 0, 0, 0, 65, 123, 206, 108, 212, 145, 69},
			},
			want: SystemEvent{
				StockLocate:    0,
				TrackingNumber: 0,
				EventCode:      EVENT_END_HOURS,
				Timestamp:      endHoursTimestamp,
			},
			wantErr: false,
		},
		{
			name: "End of Messages",
			args: args{
				data: []byte{83, 0, 0, 0, 0, 65, 193, 167, 209, 211, 137, 67},
			},
			want: SystemEvent{
				StockLocate:    0,
				TrackingNumber: 0,
				EventCode:      EVENT_END_MESSAGES,
				Timestamp:      endMessagesTimestamp,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSystemEvent(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSystemEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSystemEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemEvent_Bytes(t *testing.T) {
	startMessagesTimestamp, _ := time.ParseDuration("3h2m33.404452051s")

	tests := []struct {
		name string
		e    SystemEvent
		want []byte
	}{
		{
			name: "Start of Messages",
			e: SystemEvent{
				StockLocate:    0,
				TrackingNumber: 0,
				EventCode:      EVENT_START_MESSAGES,
				Timestamp:      startMessagesTimestamp,
			},
			want: []byte{83, 0, 0, 0, 0, 9, 246, 73, 200, 12, 211, 79},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SystemEvent.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
