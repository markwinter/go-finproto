/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/quagmt/udecimal"
)

func TestParseTradeCross(t *testing.T) {
	timestamp, _ := time.ParseDuration("9h30m0.260249144s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    TradeCross
		wantErr bool
	}{
		{
			name: "SOLY",
			args: args{data: []byte{81, 28, 191, 0, 2, 31, 26, 222, 93, 6, 56, 0, 0, 0, 0, 0, 0, 9, 66, 83, 79, 76, 89, 32, 32, 32, 32, 0, 1, 212, 192, 0, 0, 0, 0, 0, 2, 7, 156, 79}},
			want: TradeCross{
				StockLocate:    7359,
				TrackingNumber: 2,
				Timestamp:      timestamp,
				MatchNumber:    133020,
				Stock:          "SOLY",
				Shares:         2370,
				CrossPrice:     udecimal.MustParse("12"),
				CrossType:      CROSS_TYPE_NASDAQ_OPEN,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTradeCross(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTradeCross() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestTradeCross_Bytes(t *testing.T) {
	timestamp, _ := time.ParseDuration("9h30m0.260249144s")

	tests := []struct {
		name string
		tr   TradeCross
		want []byte
	}{
		{
			name: "SOLY",
			tr: TradeCross{
				StockLocate:    7359,
				TrackingNumber: 2,
				Timestamp:      timestamp,
				MatchNumber:    133020,
				Stock:          "SOLY",
				Shares:         2370,
				CrossPrice:     udecimal.MustParse("12"),
				CrossType:      CROSS_TYPE_NASDAQ_OPEN,
			},
			want: []byte{81, 28, 191, 0, 2, 31, 26, 222, 93, 6, 56, 0, 0, 0, 0, 0, 0, 9, 66, 83, 79, 76, 89, 32, 32, 32, 32, 0, 1, 212, 192, 0, 0, 0, 0, 0, 2, 7, 156, 79},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TradeCross.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
