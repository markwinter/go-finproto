/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseStockTradingAction(t *testing.T) {
	timestamp, _ := time.ParseDuration("3h7m15.46964322s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    StockTradingAction
		wantErr bool
	}{
		{
			name: "EVY",
			args: args{data: []byte{72, 10, 86, 0, 0, 10, 55, 246, 45, 77, 212, 69, 86, 89, 32, 32, 32, 32, 32, 84, 32, 32, 32, 32, 32}},
			want: StockTradingAction{
				StockLocate:    2646,
				TrackingNumber: 0,
				Timestamp:      timestamp,
				Stock:          "EVY",
				TradingState:   STATE_TRADING,
				Reserved:       32,
				Reason:         "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStockTradingAction(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStockTradingAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestStockTradingAction_Bytes(t *testing.T) {
	timestamp, _ := time.ParseDuration("3h7m15.46964322s")

	tests := []struct {
		name string
		tr   StockTradingAction
		want []byte
	}{
		{
			name: "EVY",
			tr: StockTradingAction{
				StockLocate:    2646,
				TrackingNumber: 0,
				Timestamp:      timestamp,
				Stock:          "EVY",
				TradingState:   STATE_TRADING,
				Reserved:       32,
				Reason:         "",
			},
			want: []byte{72, 10, 86, 0, 0, 10, 55, 246, 45, 77, 212, 69, 86, 89, 32, 32, 32, 32, 32, 84, 32, 32, 32, 32, 32},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StockTradingAction.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
