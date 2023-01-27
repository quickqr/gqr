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

// WithFgColor sets color that is used to draw modules (ignored if gradient is set)
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

// WithGradient will use gradient to paint modules instead of foregroundColor (if set by WithFgColor or WithFgColorRGBHex)
func WithGradient(d GradientDirection, colors ...color.Color) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.gradientConfig = &GradientConfig{d, colors}
	})
}

// WithLogo embeds image at the center of the QR
func WithLogo(img image.Image) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.logo = img
	})
}

// WithSpaceAroundLogo adds empty space behind logo so it's not drawn on top of modules
func WithSpaceAroundLogo() ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.spaceAroundLogo = true
	})
}

// WithImageSize sets size of outputted image in pixels
func WithImageSize(size int) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.size = size
	})
}

// WithQuietZone set padding around QR code.
// Note: actual size of the QR code is equal to size - quietZone * 2 (padding applied on every side )
func WithQuietZone(size int) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.quietZone = size
	})
}

// WithModuleShape sets function that will draw  modules on the image
// See: shapes.SquareModuleShape, shapes.RoundedModuleShape.
func WithModuleShape(drawer shapes.ModuleDrawer) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.drawModuleFn = drawer
	})
}

// WithFinderShape sets config for drawing 3 finders in corners of QR code.
// See: shapes.SquareFinderShape, shapes.RoundedFinderShape.
func WithFinderShape(c shapes.FinderDrawConfig) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		oo.drawFinder = c
	})
}

// WithModuleGap will set gaps between modules in percents relative to dynamic module size
// (determined by quiet zone and image size)
//
// Note: gap should be in range [0; 1). Other values are ignored
func WithModuleGap(gap float64) ImageOption {
	return newFuncOption(func(oo *imageOptions) {
		if gap < 0 || gap >= 1 {
			return
		}

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
