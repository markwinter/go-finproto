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

func TestParseOrderExecuted(t *testing.T) {
	timestamp, _ := time.ParseDuration("4h4m56.250467048s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    OrderExecuted
		wantErr bool
	}{
		{
			name: "real data",
			args: args{data: []byte{69, 31, 216, 0, 2, 13, 93, 188, 201, 226, 232, 0, 0, 0, 0, 0, 3, 224, 136, 0, 0, 3, 77, 0, 0, 0, 0, 0, 0, 73, 42}},
			want: OrderExecuted{
				StockLocate:    8152,
				Reference:      254088,
				TrackingNumber: 2,
				MatchNumber:    18730,
				Shares:         845,
				Timestamp:      timestamp,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOrderExecuted(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOrderExecuted() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestOrderExecuted_Bytes(t *testing.T) {
	timestamp, _ := time.ParseDuration("4h4m56.250467048s")

	tests := []struct {
		name string
		o    OrderExecuted
		want []byte
	}{
		{
			name: "real data",
			o: OrderExecuted{
				StockLocate:    8152,
				Reference:      254088,
				TrackingNumber: 2,
				MatchNumber:    18730,
				Shares:         845,
				Timestamp:      timestamp,
			},
			want: []byte{69, 31, 216, 0, 2, 13, 93, 188, 201, 226, 232, 0, 0, 0, 0, 0, 3, 224, 136, 0, 0, 3, 77, 0, 0, 0, 0, 0, 0, 73, 42},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderExecuted.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOrderExecutedPrice(t *testing.T) {
	timestamp, _ := time.ParseDuration("9h32m10.389319714s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    OrderExecutedPrice
		wantErr bool
	}{
		{
			name: "real data",
			args: args{data: []byte{67, 33, 73, 0, 2, 31, 57, 42, 169, 16, 34, 0, 0, 0, 0, 1, 77, 152, 216, 0, 0, 1, 144, 0, 0, 0, 0, 0, 6, 170, 14, 89, 0, 2, 85, 168}},
			want: OrderExecutedPrice{
				StockLocate:    8521,
				Reference:      21862616,
				TrackingNumber: 2,
				MatchNumber:    436750,
				Shares:         400,
				Timestamp:      timestamp,
				Printable:      true,
				ExecutionPrice: udecimal.MustParse("15.3"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOrderExecutedPrice(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOrderExecutedPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestOrderExecutedPrice_Bytes(t *testing.T) {
	timestamp, _ := time.ParseDuration("9h32m10.389319714s")

	tests := []struct {
		name string
		o    OrderExecutedPrice
		want []byte
	}{
		{
			name: "real data",
			o: OrderExecutedPrice{
				StockLocate:    8521,
				Reference:      21862616,
				TrackingNumber: 2,
				MatchNumber:    436750,
				Shares:         400,
				Timestamp:      timestamp,
				Printable:      true,
				ExecutionPrice: udecimal.MustParse("15.3"),
			},
			want: []byte{67, 33, 73, 0, 2, 31, 57, 42, 169, 16, 34, 0, 0, 0, 0, 1, 77, 152, 216, 0, 0, 1, 144, 0, 0, 0, 0, 0, 6, 170, 14, 89, 0, 2, 85, 168},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderExecutedPrice.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
