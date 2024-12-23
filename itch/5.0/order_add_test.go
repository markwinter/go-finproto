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

func TestParseOrderAdd(t *testing.T) {
	et, _ := time.ParseDuration("4h0m0.037385465s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    OrderAdd
		wantErr bool
	}{
		{
			name: "ERIC Buy",
			args: args{
				data: []byte{65, 9, 224, 0, 0, 13, 24, 197, 28, 244, 249, 0, 0, 0, 0, 0, 0, 53, 94, 66, 0, 0, 23, 112, 69, 82, 73, 67, 32, 32, 32, 32, 0, 1, 53, 196},
			},
			want: OrderAdd{
				StockLocate:    2528,
				TrackingNumber: 0,
				Stock:          "ERIC",
				Shares:         6000,
				Price:          udecimal.MustParse("7.93"),
				OrderIndicator: ORDER_INDICATOR_BUY,
				Reference:      13662,
				Timestamp:      et,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOrderAdd(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOrderAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestOrderAdd_Bytes(t *testing.T) {
	et, _ := time.ParseDuration("4h0m0.037385465s")

	tests := []struct {
		name string
		o    OrderAdd
		want []byte
	}{
		{
			name: "ERIC Buy",
			o: OrderAdd{
				StockLocate:    2528,
				TrackingNumber: 0,
				Stock:          "ERIC",
				Shares:         6000,
				Price:          udecimal.MustParse("7.93"),
				OrderIndicator: ORDER_INDICATOR_BUY,
				Reference:      13662,
				Timestamp:      et,
			},
			want: []byte{65, 9, 224, 0, 0, 13, 24, 197, 28, 244, 249, 0, 0, 0, 0, 0, 0, 53, 94, 66, 0, 0, 23, 112, 69, 82, 73, 67, 32, 32, 32, 32, 0, 1, 53, 196},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderAdd.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOrderAddAttributed(t *testing.T) {
	ct, _ := time.ParseDuration("7h43m17.935586175s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    OrderAddAttributed
		wantErr bool
	}{
		{
			name: "CODA",
			args: args{
				data: []byte{70, 6, 65, 0, 0, 25, 72, 54, 19, 123, 127, 0, 0, 0, 0, 0, 55, 179, 21, 83, 0, 0, 0, 100, 67, 79, 68, 65, 32, 32, 32, 32, 119, 53, 147, 156, 78, 73, 84, 69},
			},
			want: OrderAddAttributed{
				StockLocate:    1601,
				Reference:      3650325,
				Stock:          "CODA",
				Shares:         100,
				Price:          udecimal.MustParse("199999.99"),
				OrderIndicator: ORDER_INDICATOR_SELL,
				Timestamp:      ct,
				Attribution:    "NITE",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOrderAddAttributed(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOrderAddAttributed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestOrderAddAttributed_Bytes(t *testing.T) {
	ct, _ := time.ParseDuration("7h43m17.935586175s")

	tests := []struct {
		name string
		o    OrderAddAttributed
		want []byte
	}{
		{
			name: "CODA",
			o: OrderAddAttributed{
				StockLocate:    1601,
				Reference:      3650325,
				Stock:          "CODA",
				Shares:         100,
				Price:          udecimal.MustParse("199999.99"),
				OrderIndicator: ORDER_INDICATOR_SELL,
				Timestamp:      ct,
				Attribution:    "NITE",
			},
			want: []byte{70, 6, 65, 0, 0, 25, 72, 54, 19, 123, 127, 0, 0, 0, 0, 0, 55, 179, 21, 83, 0, 0, 0, 100, 67, 79, 68, 65, 32, 32, 32, 32, 119, 53, 147, 156, 78, 73, 84, 69},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderAddAttributed.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
