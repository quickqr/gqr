package image

import (
	"image/color"
	"reflect"
	"testing"
)

func Test_hexToRGBA(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want color.RGBA
	}{
		{
			name: "case 1",
			args: args{s: "#112233"},
			want: color.RGBA{R: 17, G: 34, B: 51, A: 255},
		},
		{
			name: "case 2",
			args: args{s: "#112"},
			want: color.RGBA{R: 17, G: 17, B: 34, A: 255},
		},
		//{
		//	name: "case 3",
		//	args: args{s: "#1122331"},
		//	want: color.RGBA{},
		//}, // panic
		//{
		//	name: "case 4",
		//	args: args{s: "#11"},
		//	want: color.RGBA{},
		//}, // panic
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFromHex(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFromHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFromColor(t *testing.T) {
	type args struct {
		c color.Color
	}
	tests := []struct {
		name string
		args args
		want color.RGBA
	}{
		{
			name: "case 0",
			args: args{
				c: color.RGBA{R: 17, G: 34, B: 51, A: 255},
			},
			want: color.RGBA{R: 17, G: 34, B: 51, A: 255},
		},
		{
			name: "case 1",
			args: args{
				c: color.Gray16{Y: 17},
			},
			want: color.RGBA{R: 17, G: 17, B: 17, A: 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFromColor(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFromColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
