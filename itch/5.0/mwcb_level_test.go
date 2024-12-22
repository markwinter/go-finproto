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

func TestParseMwcbLevel(t *testing.T) {
	timestamp, _ := time.ParseDuration("7h0m4.049455143s")

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    MwcbLevel
		wantErr bool
	}{
		{
			name: "real mwcb",
			args: args{
				data: []byte{86, 0, 0, 0, 0, 22, 236, 70, 106, 40, 39, 0, 0, 0, 70, 225, 52, 30, 128, 0, 0, 0, 66, 78, 130, 62, 64, 0, 0, 0, 60, 248, 201, 156, 0},
			},
			want: MwcbLevel{
				StockLocate:    0,
				TrackingNumber: 0,
				LevelOne:       udecimal.MustParse("3044.26"),
				LevelTwo:       udecimal.MustParse("2847.85"),
				LevelThree:     udecimal.MustParse("2618.72"),
				Timestamp:      timestamp,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMwcbLevel(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMwcbLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("%v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestMwcbLevel_Bytes(t *testing.T) {
	timestamp, _ := time.ParseDuration("7h0m4.049455143s")

	tests := []struct {
		name string
		m    MwcbLevel
		want []byte
	}{
		{
			name: "real mwcb",
			m: MwcbLevel{
				StockLocate:    0,
				TrackingNumber: 0,
				LevelOne:       udecimal.MustParse("3044.26"),
				LevelTwo:       udecimal.MustParse("2847.85"),
				LevelThree:     udecimal.MustParse("2618.72"),
				Timestamp:      timestamp,
			},
			want: []byte{86, 0, 0, 0, 0, 22, 236, 70, 106, 40, 39, 0, 0, 0, 70, 225, 52, 30, 128, 0, 0, 0, 66, 78, 130, 62, 64, 0, 0, 0, 60, 248, 201, 156, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MwcbLevel.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
