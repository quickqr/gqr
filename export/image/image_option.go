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

var DefaultImageOptions = imageOptions{
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
	logo            image.Image
	spaceAroundLogo bool

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
	color_WHITE = ParseFromHex("#ffffff")
	color_BLACK = ParseFromHex("#000000")
)

var (
	// _STATE_MAPPING mapping matrix.State to color.RGBA in debug mode.
	_STATE_MAPPING = map[gqr.QRType]color.RGBA{
		gqr.QRType_INIT:     ParseFromHex("#ffffff"), // [bg]
		gqr.QRType_DATA:     ParseFromHex("#cdc9c3"), // [bg]
		gqr.QRType_VERSION:  ParseFromHex("#000000"), // [fg]
		gqr.QRType_FORMAT:   ParseFromHex("#444444"), // [fg]
		gqr.QRType_FINDER:   ParseFromHex("#555555"), // [fg]
		gqr.QRType_DARK:     ParseFromHex("#2BA859"), // [fg]
		gqr.QRType_SPLITTER: ParseFromHex("#2BA859"), // [fg]
		gqr.QRType_TIMING:   ParseFromHex("#000000"), // [fg]
	}
)

// ParseFromHex convert hex string into color.RGBA
func ParseFromHex(s string) color.RGBA {
	c := color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0xff,
	}

	var err error
	switch len(s) {
	case 9:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 9, 7 or 4")
	}
	if err != nil {
		panic(err)
	}

	return c
}

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
