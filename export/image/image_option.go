package image

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/quickqr/gqr"
	"github.com/quickqr/gqr/export/image/shapes"
	"image"
	"image/color"
)

type ImageOption interface {
	apply(o *imageOptions)
}

var defaultImageOptions = imageOptions{
	backgroundColor: color_WHITE, // white
	foregroundColor: color_BLACK, // black
	gradientConfig:  nil,
	logo:            nil,
	size:            512,
	quietZone:       30,
	moduleGap:       0,
	drawModuleFn:    shapes.SquareModuleShape(),
	drawFinder:      shapes.SquareFinderShape(),
}

type GradientDirection = int

const (
	// GradientDirectionTLBR - Top Left -> Bottom Right
	GradientDirectionTLBR GradientDirection = iota
	// GradientDirectionTRBL - Top Right -> Bottom Left
	GradientDirectionTRBL GradientDirection = iota
)

type GradientConfig struct {
	direction GradientDirection
	colors    []color.Color
}

func dirToGradient(d GradientDirection, w float64, h float64) gg.Gradient {
	switch d {
	case GradientDirectionTLBR:
		return gg.NewLinearGradient(0, 0, w, h)
	case GradientDirectionTRBL:
		return gg.NewLinearGradient(w, 0, 0, h)
	}

	return nil
}

// imageOptions to output QR code image
type imageOptions struct {
	// backgroundColor is the background color of the QR code image.
	backgroundColor color.RGBA `default:"color.RGB{0, 0, 0}"`

	// foregroundColor is the foreground color of the QR code.
	foregroundColor color.RGBA `default:"color.RGB{1,1,1}"`

	gradientConfig *GradientConfig

	// logo this icon image would be put the center of QR Code image
	// TODO: Force color for container?
	logo                      image.Image
	drawLogoContainer         bool
	logoContainerBorderRadius float64

	// size in pixel of output image
	// Note: Actual size of the QR code will be equal to size - quietZone
	size int

	// quietZone is the size in pixels of the quiet zone around the QR code
	quietZone int

	moduleGap float64

	drawModuleFn shapes.ModuleDrawer
	drawFinder   shapes.FinderDrawConfig

	// TODO
	// halftoneImg is the halftone image for the output image.
	//halftoneImg image.Image
}

var (
	color_WHITE = parseFromHex("#ffffff")
	color_BLACK = parseFromHex("#000000")
)

var (
	// _STATE_MAPPING mapping matrix.State to color.RGBA in debug mode.
	_STATE_MAPPING = map[gqr.QRType]color.RGBA{
		gqr.QRType_INIT:     parseFromHex("#ffffff"), // [bg]
		gqr.QRType_DATA:     parseFromHex("#cdc9c3"), // [bg]
		gqr.QRType_VERSION:  parseFromHex("#000000"), // [fg]
		gqr.QRType_FORMAT:   parseFromHex("#444444"), // [fg]
		gqr.QRType_FINDER:   parseFromHex("#555555"), // [fg]
		gqr.QRType_DARK:     parseFromHex("#2BA859"), // [fg]
		gqr.QRType_SPLITTER: parseFromHex("#2BA859"), // [fg]
		gqr.QRType_TIMING:   parseFromHex("#000000"), // [fg]
	}
)

// TODO: Add support for length of 3, 6 and 8 (with alpha)

// parseFromHex convert hex string into color.RGBA
func parseFromHex(s string) color.RGBA {
	c := color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0xff,
	}

	var err error
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}
	if err != nil {
		panic(err)
	}

	return c
}

// TODO: Transparent color does not work.
func parseFromColor(c color.Color) color.RGBA {
	rgba, ok := c.(color.RGBA)
	if ok {
		return rgba
	}

	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}
