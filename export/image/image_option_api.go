package image

import (
	"image"
	"image/color"
)

// funcOption wraps a function that modifies outputImageOptions into an
// implementation of the ImageOption interface.
type funcOption struct {
	f func(oo *outputImageOptions)
}

func (fo *funcOption) apply(oo *outputImageOptions) {
	fo.f(oo)
}

func newFuncOption(f func(oo *outputImageOptions)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithBgColor background color
func WithBgColor(c color.Color) ImageOption {
	return newFuncOption(func(oo *outputImageOptions) {
		if c == nil {
			return
		}

		oo.backgroundColor = parseFromColor(c)
	})
}

// WithBgColorRGBHex background color
func WithBgColorRGBHex(hex string) ImageOption {
	return newFuncOption(func(oo *outputImageOptions) {
		if hex == "" {
			return
		}

		oo.backgroundColor = parseFromHex(hex)
	})
}

// WithFgColor QR color
func WithFgColor(c color.Color) ImageOption {
	return newFuncOption(func(oo *outputImageOptions) {
		if c == nil {
			return
		}

		oo.foregroundColor = parseFromColor(c)
	})
}

// WithFgColorRGBHex Hex string to set QR Color
func WithFgColorRGBHex(hex string) ImageOption {
	return newFuncOption(func(oo *outputImageOptions) {
		oo.foregroundColor = parseFromHex(hex)
	})
}

func WithLogo(img image.Image) ImageOption {
	return newFuncOption(func(oo *outputImageOptions) {
		oo.logo = img
	})
}

// WithImageSize specify width of each qr block
func WithImageSize(size int) ImageOption {
	return newFuncOption(func(oo *outputImageOptions) {
		oo.size = size
	})
}

// WithQuietZone set padding around the QR code.
// Note actual size of the QR code is equal to size - quietZone * 2 (2 sides)
func WithQuietZone(size int) ImageOption {
	return newFuncOption(func(oo *outputImageOptions) {
		oo.quietZone = size
	})
}

// TODO:
//// WithHalftone ...
//func WithHalftone(path string) ImageOption {
//	return newFuncOption(func(oo *outputImageOptions) {
//		srcImg, err := imgkit.Read(path)
//		if err != nil {
//			fmt.Println("Read halftone image failed: ", err)
//			return
//		}
//
//		oo.halftoneImg = srcImg
//	})
//}
