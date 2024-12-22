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

func TestParseRegSho(t *testing.T) {
	elatTimestamp, _ := time.ParseDuration("3h45m2.67061548s")
	zomTimestamp, _ := time.ParseDuration("3h7m15.847088211s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    RegSho
		wantErr bool
	}{
		{
			name: "ELAT real data",
			args: args{
				data: []byte{89, 9, 92, 0, 0, 12, 71, 213, 226, 179, 184, 69, 76, 65, 84, 32, 32, 32, 32, 48},
			},
			want: RegSho{
				Stock:          "ELAT",
				StockLocate:    2396,
				TrackingNumber: 0,
				Action:         REGSHO_NO_PRICE_TEST,
				Timestamp:      elatTimestamp,
			},
			wantErr: false,
		},
		{
			name: "ZOM real data",
			args: args{
				data: []byte{89, 34, 190, 0, 0, 10, 56, 12, 172, 168, 83, 90, 79, 77, 32, 32, 32, 32, 32, 50},
			},
			want: RegSho{
				Stock:          "ZOM",
				StockLocate:    8894,
				TrackingNumber: 0,
				Action:         REGSHO_REMAINS,
				Timestamp:      zomTimestamp,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRegSho(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRegSho() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestRegSho_Bytes(t *testing.T) {
	elatTimestamp, _ := time.ParseDuration("3h45m2.67061548s")

	tests := []struct {
		name string
		r    RegSho
		want []byte
	}{
		{
			name: "ELAT",
			r: RegSho{
				Stock:          "ELAT",
				StockLocate:    2396,
				TrackingNumber: 0,
				Action:         REGSHO_NO_PRICE_TEST,
				Timestamp:      elatTimestamp,
			},
			want: []byte{89, 9, 92, 0, 0, 12, 71, 213, 226, 179, 184, 69, 76, 65, 84, 32, 32, 32, 32, 48},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegSho.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
