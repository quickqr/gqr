package reedsolomon

import (
	"testing"
)

// table ref to: https://www.thonky.com/qr-code-tutorial/log-antilog-table
func Test_initTables(t *testing.T) {

	initTables()

	tests := []struct {
		name    string
		expIdx  int
		expWant byte
		logIdx  int
		logWant byte
	}{
		{
			name:    "case 0",
			expIdx:  0,
			expWant: 1,
			logIdx:  1,
			logWant: 0,
		},
		{
			name:    "case 1",
			expIdx:  255,
			expWant: 1,
			logIdx:  255,
			logWant: 175,
		},
		{
			name:    "case 2",
			expIdx:  300,
			expWant: gfExp[300%255],
			logIdx:  246,
			logWant: 173,
		},
		{
			name:    "case 3",
			expIdx:  380,
			expWant: 51,
			logIdx:  246,
			logWant: 173,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if v := gfExp[tt.expIdx]; v != tt.expWant {
				t.Errorf("gfExp[%d]: %v, want: %v", tt.expIdx, v, tt.expWant)
			}
			if v := gfLog[tt.logIdx]; v != tt.logWant {
				t.Errorf("gfLog[%d]: %v, want: %v", tt.logIdx, v, tt.logWant)
			}
		})
	}
}

func Test_gfMul(t *testing.T) {
	initTables()

	type args struct {
		x byte
		y byte
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{
			name: "case 0",
			args: args{
				x: 98,
				y: 7,
			},
			want: 51,
		},
		{
			name: "case 1",
			args: args{
				x: 1,
				y: 98,
			},
			want: 98,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("log x: %v, y: %v", gfLog[tt.args.x], gfLog[tt.args.y])
			// expIdx := uint(gfLog[tt.args.x]) + uint(gfLog[tt.args.y])
			// t.Logf("exp[%d]: %v", expIdx, gfExp[expIdx])
			if got := gfMul(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("gfMul() = %v, want %v", got, tt.want)
			}
		})
	}
}
