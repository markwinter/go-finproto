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

func TestParseOrderCancel(t *testing.T) {
	timestamp, _ := time.ParseDuration("4h0m9.485740746s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    OrderCancel
		wantErr bool
	}{
		{
			name: "Real data",
			args: args{
				data: []byte{88, 29, 44, 0, 0, 13, 26, 248, 71, 106, 202, 0, 0, 0, 0, 0, 0, 147, 48, 0, 0, 17, 48},
			},
			want: OrderCancel{
				StockLocate: 7468,
				Reference:   37680,
				Shares:      4400,
				Timestamp:   timestamp,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOrderCancel(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOrderCancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestOrderCancel_Bytes(t *testing.T) {
	timestamp, _ := time.ParseDuration("4h0m9.485740746s")

	tests := []struct {
		name string
		o    OrderCancel
		want []byte
	}{
		{
			name: "real data",
			o: OrderCancel{
				StockLocate: 7468,
				Reference:   37680,
				Shares:      4400,
				Timestamp:   timestamp,
			},
			want: []byte{88, 29, 44, 0, 0, 13, 26, 248, 71, 106, 202, 0, 0, 0, 0, 0, 0, 147, 48, 0, 0, 17, 48},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderCancel.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
