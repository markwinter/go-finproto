/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseStockDirectory(t *testing.T) {
	timestamp, _ := time.ParseDuration("3h7m14.939262551s")
	bTimestamp, _ := time.ParseDuration("3h7m14.938300595s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    StockDirectory
		wantErr bool
	}{
		{
			name: "",
			args: args{data: []byte{82, 3, 81, 0, 0, 10, 55, 214, 144, 86, 87, 66, 72, 66, 32, 32, 32, 32, 32, 65, 32, 0, 0, 0, 100, 78, 67, 90, 32, 80, 78, 32, 50, 78, 0, 0, 0, 0, 78}},
			want: StockDirectory{
				Timestamp:                   timestamp,
				Stock:                       "BHB",
				StockLocate:                 849,
				TrackingNumber:              0,
				RoundLotSize:                100,
				RoundLotsOnly:               false,
				IssueSubType:                ICS_NOT_APPLICABLE,
				IssueClassification:         IC_COMMON_STOCK,
				InverseIndicator:            false,
				Authenticity:                AUTHENTICITY_LIVE,
				EtpLeverageFactor:           0,
				ShortSaleThresholdIndicator: "N",
				IpoFlag:                     " ",
				LuldReferencePriceTier:      "2",
				EtpFlag:                     "N",
				MarketCategory:              MKTCTG_NYSE_AMERICAN,
				FinancialStatusIndicator:    FSI_NOT_AVAILABLE,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{data: []byte{82, 3, 45, 0, 0, 10, 55, 214, 129, 168, 179, 66, 69, 83, 84, 32, 32, 32, 32, 78, 32, 0, 0, 0, 100, 78, 65, 90, 32, 80, 78, 32, 50, 78, 0, 0, 0, 1, 78}},
			want: StockDirectory{
				Timestamp:                   bTimestamp,
				Stock:                       "BEST",
				StockLocate:                 813,
				TrackingNumber:              0,
				RoundLotSize:                100,
				RoundLotsOnly:               false,
				IssueSubType:                ICS_NOT_APPLICABLE,
				IssueClassification:         IC_AMERICAN_DEPOSITORY_SHARE,
				InverseIndicator:            false,
				Authenticity:                AUTHENTICITY_LIVE,
				EtpLeverageFactor:           1,
				ShortSaleThresholdIndicator: "N",
				IpoFlag:                     " ",
				LuldReferencePriceTier:      "2",
				EtpFlag:                     "N",
				MarketCategory:              MKTCTG_NYSE,
				FinancialStatusIndicator:    FSI_NOT_AVAILABLE,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStockDirectory(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStockDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestStockDirectory_Bytes(t *testing.T) {
	timestamp, _ := time.ParseDuration("3h7m14.939262551s")
	bTimestamp, _ := time.ParseDuration("3h7m14.938300595s")

	tests := []struct {
		name string
		s    StockDirectory
		want []byte
	}{
		{
			name: "BHB",
			s: StockDirectory{
				Timestamp:                   timestamp,
				Stock:                       "BHB",
				StockLocate:                 849,
				TrackingNumber:              0,
				RoundLotSize:                100,
				RoundLotsOnly:               false,
				IssueSubType:                ICS_NOT_APPLICABLE,
				IssueClassification:         IC_COMMON_STOCK,
				InverseIndicator:            false,
				Authenticity:                AUTHENTICITY_LIVE,
				EtpLeverageFactor:           0,
				ShortSaleThresholdIndicator: "N",
				IpoFlag:                     " ",
				LuldReferencePriceTier:      "2",
				EtpFlag:                     "N",
				MarketCategory:              MKTCTG_NYSE_AMERICAN,
				FinancialStatusIndicator:    FSI_NOT_AVAILABLE,
			},
			want: []byte{82, 3, 81, 0, 0, 10, 55, 214, 144, 86, 87, 66, 72, 66, 32, 32, 32, 32, 32, 65, 32, 0, 0, 0, 100, 78, 67, 90, 32, 80, 78, 32, 50, 78, 0, 0, 0, 0, 78},
		},
		{
			name: "BHB",
			s: StockDirectory{
				Timestamp:                   bTimestamp,
				Stock:                       "BEST",
				StockLocate:                 813,
				TrackingNumber:              0,
				RoundLotSize:                100,
				RoundLotsOnly:               false,
				IssueSubType:                ICS_NOT_APPLICABLE,
				IssueClassification:         IC_AMERICAN_DEPOSITORY_SHARE,
				InverseIndicator:            false,
				Authenticity:                AUTHENTICITY_LIVE,
				EtpLeverageFactor:           1,
				ShortSaleThresholdIndicator: "N",
				IpoFlag:                     " ",
				LuldReferencePriceTier:      "2",
				EtpFlag:                     "N",
				MarketCategory:              MKTCTG_NYSE,
				FinancialStatusIndicator:    FSI_NOT_AVAILABLE,
			},
			want: []byte{82, 3, 45, 0, 0, 10, 55, 214, 129, 168, 179, 66, 69, 83, 84, 32, 32, 32, 32, 78, 32, 0, 0, 0, 100, 78, 65, 90, 32, 80, 78, 32, 50, 78, 0, 0, 0, 1, 78},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Bytes(); !cmp.Equal(got, tt.want) {
				t.Errorf("StockDirectory.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
