package image

import (
	"github.com/quickqr/gqr/export/image/shapes"
	"image"
	"image/color"
)

// funcOption wraps a function that modifies imageOptions into an
// implementation of the ImageOption interface.
type funcOption struct {
	f func(oo *imageOptions)
}

func (fo *funcOption) apply(oo *imageOptions) {
	fo.f(oo)
}

func newFuncOption(f func(oo *imageOptions)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithBgColor background color
func WithBgColor(c color.Color) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		if c == nil {
			return
		}

		oo.backgroundColor = parseFromColor(c)
	})
}

// WithBgColorRGBHex background color
func WithBgColorRGBHex(hex string) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		if hex == "" {
			return
		}

		oo.backgroundColor = parseFromHex(hex)
	})
}

// WithFgColor QR color
func WithFgColor(c color.Color) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		if c == nil {
			return
		}

		oo.foregroundColor = parseFromColor(c)
	})
}

// WithFgColorRGBHex Hex string to set QR Color
func WithFgColorRGBHex(hex string) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.foregroundColor = parseFromHex(hex)
	})
}

func WithLogo(img image.Image) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.logo = img
	})
}

// WithImageSize specify width of each qr block
func WithImageSize(size int) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.size = size
	})
}

// WithQuietZone set padding around the QR code.
// Note actual size of the QR code is equal to size - quietZone * 2 (2 sides)
func WithQuietZone(size int) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.quietZone = size
	})
}

// WithModuleShape sets function that will draw  modules on the image
func WithModuleShape(drawer shapes.ModuleShapeDrawer) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.drawModuleFn = drawer
	})
}

// WithModuleGap sets margin between each module on the QR code
func WithModuleGap(gap float64) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.moduleGap = gap
	})
}

// TODO:
//// WithHalftone ...
//func WithHalftone(path string) ImageOption {
//	return newFuncOption(func(oo *imageOptions) {
//		srcImg, err := imgkit.Read(path)
//		if err != nil {
//			fmt.Println("Read halftone image failed: ", err)
//			return
//		}
//
//		oo.halftoneImg = srcImg
//	})
//}
