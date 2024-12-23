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

func TestParseOrderReplace(t *testing.T) {
	timestamp, _ := time.ParseDuration("4h0m27.588626843s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    OrderReplace
		wantErr bool
	}{
		{
			name: "real data",
			args: args{data: []byte{85, 9, 212, 0, 0, 13, 31, 47, 75, 137, 155, 0, 0, 0, 0, 0, 1, 24, 254, 0, 0, 0, 0, 0, 1, 39, 34, 0, 0, 11, 184, 0, 2, 208, 180}},
			want: OrderReplace{
				StockLocate:       2516,
				TrackingNumber:    0,
				OriginalReference: 71934,
				NewReference:      75554,
				Shares:            3000,
				Price:             udecimal.MustParse("18.45"),
				Timestamp:         timestamp,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOrderReplace(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOrderReplace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestOrderReplace_Bytes(t *testing.T) {
	timestamp, _ := time.ParseDuration("4h0m27.588626843s")

	tests := []struct {
		name string
		o    OrderReplace
		want []byte
	}{
		{
			name: "real data",
			o: OrderReplace{
				StockLocate:       2516,
				TrackingNumber:    0,
				OriginalReference: 71934,
				NewReference:      75554,
				Shares:            3000,
				Price:             udecimal.MustParse("18.45"),
				Timestamp:         timestamp,
			},
			want: []byte{85, 9, 212, 0, 0, 13, 31, 47, 75, 137, 155, 0, 0, 0, 0, 0, 1, 24, 254, 0, 0, 0, 0, 0, 1, 39, 34, 0, 0, 11, 184, 0, 2, 208, 180},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderReplace.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
