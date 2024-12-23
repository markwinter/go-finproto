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

func TestParseOrderDelete(t *testing.T) {
	timestamp, _ := time.ParseDuration("4h0m4.757910745s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    OrderDelete
		wantErr bool
	}{
		{
			name: "real data",
			args: args{data: []byte{68, 31, 163, 0, 0, 13, 25, 222, 122, 116, 217, 0, 0, 0, 0, 0, 0, 5, 208}},
			want: OrderDelete{
				StockLocate:    8099,
				TrackingNumber: 0,
				Reference:      1488,
				Timestamp:      timestamp,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOrderDelete(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOrderDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestOrderDelete_Bytes(t *testing.T) {
	timestamp, _ := time.ParseDuration("4h0m4.757910745s")

	tests := []struct {
		name string
		o    OrderDelete
		want []byte
	}{
		{
			name: "real data",
			o: OrderDelete{
				StockLocate:    8099,
				TrackingNumber: 0,
				Reference:      1488,
				Timestamp:      timestamp,
			},
			want: []byte{68, 31, 163, 0, 0, 13, 25, 222, 122, 116, 217, 0, 0, 0, 0, 0, 0, 5, 208},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderDelete.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
