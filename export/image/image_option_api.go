package image

import (
	"github.com/quickqr/gqr/export/image/shapes"
	"image"
	"image/color"
)

// funcOption wraps a function that modifies exportOptions into an
// implementation of the ExportOption interface.
type funcOption struct {
	f func(oo *exportOptions)
}

func (fo *funcOption) apply(oo *exportOptions) {
	fo.f(oo)
}

func newFuncOption(f func(oo *exportOptions)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithBgColor background color
func WithBgColor(c color.Color) ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		if c == nil {
			return
		}

		oo.backgroundColor = parseFromColor(c)
	})
}

// WithBgColorHex background color
func WithBgColorHex(hex string) ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		if hex == "" {
			return
		}

		oo.backgroundColor = ParseFromHex(hex)
	})
}

// WithFgColor sets color that is used to draw modules (ignored if gradient is set)
func WithFgColor(c color.Color) ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		if c == nil {
			return
		}

		oo.foregroundColor = parseFromColor(c)
	})
}

// WithFgColorHex Hex string to set QR Color
func WithFgColorHex(hex string) ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		oo.foregroundColor = ParseFromHex(hex)
	})
}

// WithGradient will use gradient to paint modules instead of foregroundColor (if set by WithFgColor or WithFgColorHex)
func WithGradient(d GradientDirection, colors ...color.Color) ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		oo.gradientConfig = &GradientConfig{d, colors}
	})
}

// WithLogo embeds image at the center of the QR
func WithLogo(img image.Image) ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		oo.logo = img
	})
}

// WithSpaceAroundLogo adds empty space behind logo, so it's not drawn on top of modules
func WithSpaceAroundLogo() ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		oo.spaceAroundLogo = true
	})
}

// WithImageSize sets size of outputted image in pixels
// Values less  than 1 are ignored
func WithImageSize(size int) ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		if size < 1 {
			return
		}

		oo.size = size
	})
}

// WithQuietZone set padding around QR code.
// Note: actual size of the QR code is equal to size - quietZone * 2 (padding applied on every side )
func WithQuietZone(size int) ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		oo.quietZone = size
	})
}

// WithModuleShape sets function that will draw  modules on the image
// See: shapes.SquareModuleShape, shapes.RoundedModuleShape.
func WithModuleShape(drawer shapes.ModuleDrawer) ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		oo.moduleDrawer = drawer
	})
}

// WithFinderShape sets config for drawing 3 finders in corners of QR code.
// See: shapes.SquareFinderShape, shapes.RoundedFinderShape.
func WithFinderShape(c shapes.FinderDrawConfig) ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		oo.finderDrawer = c
	})
}

// WithModuleGap will set gaps between modules in percents relative to dynamic module size
// (determined by quiet zone and image size)
//
// Note: gap should be in range [0; 1). Other values are ignored
func WithModuleGap(gap float64) ExportOption {
	return newFuncOption(func(oo *exportOptions) {
		if gap < 0 || gap >= 1 {
			return
		}

		oo.moduleGap = gap
	})
}

// TODO:
//// WithHalftone ...
//func WithHalftone(path string) ExportOption {
//	return newFuncOption(func(oo *exportOptions) {
//		srcImg, err := imgkit.Read(path)
//		if err != nil {
//			fmt.Println("Read halftone image failed: ", err)
//			return
//		}
//
//		oo.halftoneImg = srcImg
//	})
//}
