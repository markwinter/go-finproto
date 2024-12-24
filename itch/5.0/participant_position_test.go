/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseParticipantPosition(t *testing.T) {
	nmrkT, _ := time.ParseDuration("3h7m25.913251167s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    ParticipantPosition
		wantErr bool
	}{
		{
			name: "NMRK",
			args: args{data: []byte{76, 21, 203, 0, 0, 10, 58, 100, 170, 29, 95, 67, 79, 87, 78, 78, 77, 82, 75, 32, 32, 32, 32, 89, 78, 65}},
			want: ParticipantPosition{
				StockLocate:    5579,
				Mpid:           "COWN",
				Stock:          "NMRK",
				TrackingNumber: 0,
				PrimaryMM:      true,
				Mode:           MMMODE_NORMAL,
				State:          MMSTATE_ACTIVE,
				Timestamp:      nmrkT,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseParticipantPosition(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseParticipantPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestParticipantPosition_Bytes(t *testing.T) {
	nmrkT, _ := time.ParseDuration("3h7m25.913251167s")

	tests := []struct {
		name string
		p    ParticipantPosition
		want []byte
	}{
		{
			name: "NMRK",
			p: ParticipantPosition{
				StockLocate:    5579,
				Mpid:           "COWN",
				Stock:          "NMRK",
				TrackingNumber: 0,
				PrimaryMM:      true,
				Mode:           MMMODE_NORMAL,
				State:          MMSTATE_ACTIVE,
				Timestamp:      nmrkT,
			},
			want: []byte{76, 21, 203, 0, 0, 10, 58, 100, 170, 29, 95, 67, 79, 87, 78, 78, 77, 82, 75, 32, 32, 32, 32, 89, 78, 65},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Bytes(); !cmp.Equal(got, tt.want) {
				t.Errorf("ParticipantPosition.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
