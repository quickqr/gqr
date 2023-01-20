// Package reedsolomon ...
// ref to doc: http://www.drdobbs.com/testing/error-correction-with-reed-solomon/240157266
// ref to project: github.com/skip2/go-qrcode/reedsolomon
package reedsolomon

import (
	"reflect"
	"testing"

	"github.com/quickqr/gqr/reedsolomon/binary"
)

// func TestEncode(t *testing.T) {
// 	/*  10 error correction
// 	x10 + α251x9 + α67x8 + α46x7 + α61x6 + α118x5 + α70x4+ α64x3 + α94x2 + α32x + α45
// 	*/
// 	bin := binary.New()

// 	bin.AppendBytes([]byte{
// 		0x40, 0xd2, 0x75, 0x47, 0x76, 0x17, 0x32, 0x06,
// 		0x27, 0x26, 0x96, 0xc6, 0xc6, 0x96, 0x70, 0xec}...)
// 	t.Logf("bin: %v\n", bin.Bytes())

// 	bout := Encode(bin, 10)
// 	// want remainder: 0xbc 0x2a 0x90 0x13 0x6b 0xaf 0xef 0xfd 0x4b 0xe0
// 	t.Logf("bout %v\n", bout.Bytes())
// }

func TestEncode(t *testing.T) {
	type args struct {
		bin        *binary.Binary
		numECWords int
		data       []byte
	}
	tests := []struct {
		name     string
		args     args
		want     *binary.Binary
		wantData []byte
	}{
		{
			name: "case 0",
			args: args{
				bin:        binary.New(),
				numECWords: 10,
				data:       []byte{32, 91, 11, 120, 209, 114, 220, 77, 67, 64, 236, 17, 236},
			},
			want:     binary.New(),
			wantData: []byte{87, 86, 68, 17, 99, 235, 189, 232, 98, 195},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.bin.AppendBytes(tt.args.data...)

			tt.want.AppendBytes(tt.args.data...)
			tt.want.AppendBytes(tt.wantData...)
			if got := Encode(tt.args.bin, tt.args.numECWords); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %3v, want %3v", got.Bytes(), tt.want.Bytes())
			}
		})
	}
}
