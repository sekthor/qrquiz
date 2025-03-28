package domain

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestBitmap_Value(t *testing.T) {
	tests := []struct {
		name    string
		b       *Bitmap
		want    driver.Value
		wantErr bool
	}{
		{
			name: "Serialize Bitmap",
			b: &Bitmap{
				{true, false, true, false, true},
				{false, true, false, true, false},
				{true, false, true, false, true},
				{false, true, false, true, false},
				{true, false, true, false, true},
			},
			want: []byte{
				0x01, 0x00, 0x01, 0x00, 0x01,
				0x00, 0x01, 0x00, 0x01, 0x00,
				0x01, 0x00, 0x01, 0x00, 0x01,
				0x00, 0x01, 0x00, 0x01, 0x00,
				0x01, 0x00, 0x01, 0x00, 0x01,
			},
			wantErr: false,
		},
		{
			name: "not square: row too big",
			b: &Bitmap{
				{true, false, true, false, true},
				{false, true, false, true, false},
				{true, false, true, false, true},
				{false, true, false, true, false},
				{true, false, true, false, true, false},
			},
			want: []byte{
				0x01, 0x00, 0x01, 0x00, 0x01,
				0x00, 0x01, 0x00, 0x01, 0x00,
				0x01, 0x00, 0x01, 0x00, 0x01,
				0x00, 0x01, 0x00, 0x01, 0x00,
				0x01, 0x00, 0x01, 0x00, 0x01, 0x00,
			},
			wantErr: false,
		},
		{
			name: "not square: row too small",
			b: &Bitmap{
				{true, false, true, false, true},
				{false, true, false, true, false},
				{true, false, true, false, true},
				{false, true, false, true, false},
				{true, false, true, false},
			},
			want: []byte{
				0x01, 0x00, 0x01, 0x00, 0x01,
				0x00, 0x01, 0x00, 0x01, 0x00,
				0x01, 0x00, 0x01, 0x00, 0x01,
				0x00, 0x01, 0x00, 0x01, 0x00,
				0x01, 0x00, 0x01, 0x00,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bitmap.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bitmap.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitmap_Scan(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		b       *Bitmap
		args    args
		wantErr bool
	}{
		{
			name: "Serialize Bitmap",
			b: &Bitmap{
				{true, false, true, false, true},
				{false, true, false, true, false},
				{true, false, true, false, true},
				{false, true, false, true, false},
				{true, false, true, false, true},
			},
			args: args{
				[]byte{
					0x01, 0x00, 0x01, 0x00, 0x01,
					0x00, 0x01, 0x00, 0x01, 0x00,
					0x01, 0x00, 0x01, 0x00, 0x01,
					0x00, 0x01, 0x00, 0x01, 0x00,
					0x01, 0x00, 0x01, 0x00, 0x01,
				},
			},
			wantErr: false,
		},
		{
			name: "not square: row too big",
			args: args{
				[]byte{
					0x01, 0x00, 0x01, 0x00, 0x01,
					0x00, 0x01, 0x00, 0x01, 0x00,
					0x01, 0x00, 0x01, 0x00, 0x01,
					0x00, 0x01, 0x00, 0x01, 0x00,
					0x01, 0x00, 0x01, 0x00, 0x01, 0x00,
				},
			},
			wantErr: true,
		},
		{
			name: "not square: row too small",
			args: args{
				[]byte{
					0x01, 0x00, 0x01, 0x00, 0x01,
					0x00, 0x01, 0x00, 0x01, 0x00,
					0x01, 0x00, 0x01, 0x00, 0x01,
					0x00, 0x01, 0x00, 0x01, 0x00,
					0x01, 0x00, 0x01, 0x00,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Bitmap.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
