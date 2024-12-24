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

func TestParseTradeNonCross(t *testing.T) {
	timestamp, _ := time.ParseDuration("6h42m41.814553968s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    TradeNonCross
		wantErr bool
	}{
		{
			name: "FB",
			args: args{data: []byte{80, 10, 181, 0, 2, 21, 249, 156, 95, 165, 112, 0, 0, 0, 0, 0, 0, 0, 0, 66, 0, 0, 0, 250, 70, 66, 32, 32, 32, 32, 32, 32, 0, 31, 149, 240, 0, 0, 0, 0, 0, 0, 113, 107}},
			want: TradeNonCross{
				StockLocate:    2741,
				TrackingNumber: 2,
				MatchNumber:    29035,
				Timestamp:      timestamp,
				Stock:          "FB",
				Reference:      0,
				Shares:         250,
				Price:          udecimal.MustParse("207"),
				OrderIndicator: ORDER_INDICATOR_BUY,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTradeNonCross(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTradeNonCross() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestTradeNonCross_Bytes(t *testing.T) {
	timestamp, _ := time.ParseDuration("6h42m41.814553968s")

	tests := []struct {
		name string
		tr   TradeNonCross
		want []byte
	}{
		{
			name: "FB",
			tr: TradeNonCross{
				StockLocate:    2741,
				TrackingNumber: 2,
				MatchNumber:    29035,
				Timestamp:      timestamp,
				Stock:          "FB",
				Reference:      0,
				Shares:         250,
				Price:          udecimal.MustParse("207"),
				OrderIndicator: ORDER_INDICATOR_BUY,
			},
			want: []byte{80, 10, 181, 0, 2, 21, 249, 156, 95, 165, 112, 0, 0, 0, 0, 0, 0, 0, 0, 66, 0, 0, 0, 250, 70, 66, 32, 32, 32, 32, 32, 32, 0, 31, 149, 240, 0, 0, 0, 0, 0, 0, 113, 107},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TradeNonCross.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
