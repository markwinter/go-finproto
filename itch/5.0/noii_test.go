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

func TestParseNoii(t *testing.T) {
	xmhqTimestamp, _ := time.ParseDuration("9h28m40.088991343s")
	pruTimestamp, _ := time.ParseDuration("9h28m40.087987805s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Noii
		wantErr bool
	}{
		{
			name: "XMHQ",
			args: args{
				data: []byte{73, 34, 63, 0, 0, 31, 8, 51, 200, 182, 111, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 79, 88, 77, 72, 81, 32, 32, 32, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 79, 32},
			},
			want: Noii{
				StockLocate:        8767,
				TrackingNumber:     0,
				Timestamp:          xmhqTimestamp,
				Stock:              "XMHQ",
				PairedShares:       0,
				ImbalanceShares:    0,
				ImbalanceDirection: IMBALANCE_INSUFFICIENT,
				CrossType:          CROSS_TYPE_NASDAQ_OPEN,
				VariationIndicator: 32,
			},
			wantErr: false,
		},
		{
			name: "PRU",
			args: args{
				data: []byte{73, 24, 233, 0, 0, 31, 8, 51, 185, 102, 93, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 79, 80, 82, 85, 32, 32, 32, 32, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 79, 32},
			},
			want: Noii{
				StockLocate:        6377,
				TrackingNumber:     0,
				Timestamp:          pruTimestamp,
				Stock:              "PRU",
				PairedShares:       0,
				ImbalanceShares:    0,
				ImbalanceDirection: IMBALANCE_INSUFFICIENT,
				CrossType:          CROSS_TYPE_NASDAQ_OPEN,
				VariationIndicator: 32,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNoii(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNoii() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestNoii_Bytes(t *testing.T) {
	xmhqTimestamp, _ := time.ParseDuration("9h28m40.088991343s")
	pruTimestamp, _ := time.ParseDuration("9h28m40.087987805s")

	tests := []struct {
		name string
		n    Noii
		want []byte
	}{
		{
			name: "PRU",
			n: Noii{
				StockLocate:        6377,
				TrackingNumber:     0,
				Timestamp:          pruTimestamp,
				Stock:              "PRU",
				PairedShares:       0,
				ImbalanceShares:    0,
				ImbalanceDirection: IMBALANCE_INSUFFICIENT,
				CrossType:          CROSS_TYPE_NASDAQ_OPEN,
				VariationIndicator: 32,
			},
			want: []byte{73, 24, 233, 0, 0, 31, 8, 51, 185, 102, 93, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 79, 80, 82, 85, 32, 32, 32, 32, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 79, 32},
		},
		{
			name: "XMHQ",
			n: Noii{
				StockLocate:        8767,
				TrackingNumber:     0,
				Timestamp:          xmhqTimestamp,
				Stock:              "XMHQ",
				PairedShares:       0,
				ImbalanceShares:    0,
				ImbalanceDirection: IMBALANCE_INSUFFICIENT,
				CrossType:          CROSS_TYPE_NASDAQ_OPEN,
				VariationIndicator: 32,
			},
			want: []byte{73, 34, 63, 0, 0, 31, 8, 51, 200, 182, 111, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 79, 88, 77, 72, 81, 32, 32, 32, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 79, 32},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Noii.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
