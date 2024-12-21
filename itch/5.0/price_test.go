package itch

import (
	"reflect"
	"testing"

	"github.com/quagmt/udecimal"
)

func Test_bytesToPrice(t *testing.T) {
	type args struct {
		data      []byte
		precision uint8
	}
	tests := []struct {
		name    string
		args    args
		want    udecimal.Decimal
		wantErr bool
	}{
		{
			name: "test correct convert precision 4",
			args: args{
				data:      []byte{0x00, 0x02, 0xE6, 0x30},
				precision: 4,
			},
			want:    udecimal.MustParse("19.0000"),
			wantErr: false,
		},
		{
			name: "test correct convert precision 4",
			args: args{
				data:      []byte{0x00, 0x01, 0xD4, 0xC0},
				precision: 4,
			},
			want:    udecimal.MustParse("12.0000"),
			wantErr: false,
		},
		{
			name: "test correct convert precision 4",
			args: args{
				data:      []byte{0x04, 0x9C, 0xD6, 0x50},
				precision: 4,
			},
			want:    udecimal.MustParse("7738.7344"),
			wantErr: false,
		},
		{
			name: "test correct convert precision 4 with right padded zeroes",
			args: args{
				data:      []byte{0x20, 0x03, 0x45, 0xE4},
				precision: 4,
			},
			want:    udecimal.MustParse("21.4500"),
			wantErr: false,
		},
		{
			name: "test correct convert precision 4 with left padded spaces",
			args: args{
				data:      []byte{0x20, 0x20, 0x3A, 0x98},
				precision: 4,
			},
			want:    udecimal.MustParse("1.5000"),
			wantErr: false,
		},
		{
			name: "test correct convert precision 8 with left padded spaces",
			args: args{
				data:      []byte{0x20, 0x20, 0x20, 0x7F, 0x04, 0x99, 0x44, 0x80},
				precision: 8,
			},
			want:    udecimal.MustParse("5455.38000000"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bytesToPrice(tt.args.data, tt.args.precision)
			if (err != nil) != tt.wantErr {
				t.Errorf("bytesToPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bytesToPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_priceToBytes(t *testing.T) {
	type args struct {
		price     udecimal.Decimal
		precision uint8
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test correct convert precision 4",
			args: args{
				price:     udecimal.MustParse("7738.7344"),
				precision: 4,
			},
			want:    []byte{0x04, 0x9C, 0xD6, 0x50},
			wantErr: false,
		},
		{
			name: "test correct convert precision 4 with right padded zeroes",
			args: args{
				price:     udecimal.MustParse("21.4500"),
				precision: 4,
			},
			want:    []byte{0x20, 0x03, 0x45, 0xE4},
			wantErr: false,
		},
		{
			name: "test correct convert precision 4 with left padded spaces",
			args: args{
				price:     udecimal.MustParse("1.5000"),
				precision: 4,
			},
			want:    []byte{0x20, 0x20, 0x3A, 0x98},
			wantErr: false,
		},
		{
			name: "test correct convert precision 8 with left padded spaces",
			args: args{
				price:     udecimal.MustParse("5455.38000000"),
				precision: 8,
			},
			want:    []byte{0x20, 0x20, 0x20, 0x7F, 0x04, 0x99, 0x44, 0x80},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := priceToBytes(tt.args.price, tt.args.precision)
			if (err != nil) != tt.wantErr {
				t.Errorf("priceToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("priceToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
