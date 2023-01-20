package reedsolomon

import (
	"reflect"
	"testing"
)

func Test_rsGenPoly(t *testing.T) {
	type args struct {
		numECWords int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "case 0",
			args: args{
				numECWords: 2,
			},
			want: []byte{1, 3, 2},
		},
		{
			name: "case 0",
			args: args{
				numECWords: 3,
			},
			want: []byte{1, 7, 14, 8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rsGenPoly(tt.args.numECWords); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rsGenPoly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_polyDiv(t *testing.T) {
	type args struct {
		dividend []byte
		divisor  []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "case 0",
			args: args{
				dividend: []byte{12, 34, 56, 23},
				divisor:  rsGenPoly(3),
			},
			want: []byte{107, 77, 39},
		},
		{
			name: "case 1",
			args: args{
				dividend: []byte{32, 91, 11, 120, 209, 114, 220, 77, 67, 64, 236, 17, 236},
				divisor:  rsGenPoly(10),
			},
			want: []byte{87, 86, 68, 17, 99, 235, 189, 232, 98, 195},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("divisor %v", tt.args.divisor)
			got := polyDiv(tt.args.dividend, tt.args.divisor)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("polyDiv() got = %v, want %v", got, tt.want)
			// }
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("polyDiv() got = %v, want %v", got, tt.want)
			}
		})
	}
}
