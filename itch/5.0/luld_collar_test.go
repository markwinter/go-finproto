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

func TestParseLuldCollar(t *testing.T) {
	rkdaTimestamp, _ := time.ParseDuration("10h5m52.949417833s")
	dtssTimestamp, _ := time.ParseDuration("9h42m29.889620336s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    LuldCollar
		wantErr bool
	}{
		{
			name: "RKDA",
			args: args{
				data: []byte{74, 26, 108, 0, 0, 33, 16, 20, 162, 247, 105, 82, 75, 68, 65, 32, 32, 32, 32, 0, 0, 212, 28, 0, 0, 233, 52, 0, 0, 173, 112, 0, 0, 0, 1},
			},
			want: LuldCollar{
				StockLocate:    6764,
				TrackingNumber: 0,
				Timestamp:      rkdaTimestamp,
				Stock:          "RKDA",
				ReferencePrice: udecimal.MustParse("5.43"),
				UpperPrice:     udecimal.MustParse("5.97"),
				LowerPrice:     udecimal.MustParse("4.44"),
				Extension:      1,
			},
			wantErr: false,
		},
		{
			name: "DTSS",
			args: args{
				data: []byte{74, 8, 132, 0, 0, 31, 201, 103, 193, 121, 112, 68, 84, 83, 83, 32, 32, 32, 32, 0, 0, 154, 76, 0, 0, 162, 28, 0, 0, 66, 4, 0, 0, 0, 0},
			},
			want: LuldCollar{
				StockLocate:    2180,
				TrackingNumber: 0,
				Timestamp:      dtssTimestamp,
				Stock:          "DTSS",
				ReferencePrice: udecimal.MustParse("3.95"),
				UpperPrice:     udecimal.MustParse("4.15"),
				LowerPrice:     udecimal.MustParse("1.69"),
				Extension:      0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLuldCollar(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLuldCollar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestLuldCollar_Bytes(t *testing.T) {
	rkdaTimestamp, _ := time.ParseDuration("10h5m52.949417833s")
	dtssTimestamp, _ := time.ParseDuration("9h42m29.889620336s")

	tests := []struct {
		name string
		l    LuldCollar
		want []byte
	}{
		{
			name: "RKDA",
			l: LuldCollar{
				StockLocate:    6764,
				TrackingNumber: 0,
				Timestamp:      rkdaTimestamp,
				Stock:          "RKDA",
				ReferencePrice: udecimal.MustParse("5.43"),
				UpperPrice:     udecimal.MustParse("5.97"),
				LowerPrice:     udecimal.MustParse("4.44"),
				Extension:      1,
			},
			want: []byte{74, 26, 108, 0, 0, 33, 16, 20, 162, 247, 105, 82, 75, 68, 65, 32, 32, 32, 32, 0, 0, 212, 28, 0, 0, 233, 52, 0, 0, 173, 112, 0, 0, 0, 1},
		},
		{
			name: "DTSS",
			l: LuldCollar{
				StockLocate:    2180,
				TrackingNumber: 0,
				Timestamp:      dtssTimestamp,
				Stock:          "DTSS",
				ReferencePrice: udecimal.MustParse("3.95"),
				UpperPrice:     udecimal.MustParse("4.15"),
				LowerPrice:     udecimal.MustParse("1.69"),
				Extension:      0,
			},
			want: []byte{74, 8, 132, 0, 0, 31, 201, 103, 193, 121, 112, 68, 84, 83, 83, 32, 32, 32, 32, 0, 0, 154, 76, 0, 0, 162, 28, 0, 0, 66, 4, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LuldCollar.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
