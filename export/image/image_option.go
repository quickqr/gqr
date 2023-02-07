package image

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/quickqr/gqr/export/image/shapes"
	"image"
	"image/color"
)

type ExportOption interface {
	apply(o *exportOptions)
}

var DefaultImageOptions = exportOptions{
	backgroundColor: color_WHITE, // white
	foregroundColor: color_BLACK, // black
	gradientConfig:  nil,
	logo:            nil,
	logoScale:       0.8,
	size:            512,
	quietZone:       30,
	moduleGap:       0,
	moduleDrawer:    shapes.SquareModuleShape(),
	finderDrawer:    shapes.SquareFinderShape(),
}

type GradientDirection = int

const (
	// GradientDirectionLTR - Top Left -> Bottom Right
	GradientDirectionLTR GradientDirection = iota
	// GradientDirectionRTL - Top Right -> Bottom Left
	GradientDirectionRTL GradientDirection = iota
)

type GradientConfig struct {
	direction GradientDirection
	colors    []color.Color
}

func dirToGradient(d GradientDirection, w float64, h float64) gg.Gradient {
	switch d {
	case GradientDirectionLTR:
		return gg.NewLinearGradient(0, 0, w, h)
	case GradientDirectionRTL:
		return gg.NewLinearGradient(w, 0, 0, h)
	}

	return nil
}

// exportOptions to output QR code image
type exportOptions struct {
	// backgroundColor is the background color of the QR code image.
	backgroundColor color.RGBA

	// foregroundColor is the foreground color of the QR code.
	foregroundColor color.RGBA

	gradientConfig *GradientConfig

	// logo this icon image would be put the center of QR Code image
	// TODO: Force color for container?
	logo            image.Image
	logoScale       float64
	spaceAroundLogo bool

	// size in pixel of output image
	// Note: Actual size of the QR code will be equal to size - quietZone
	size int

	// quietZone is the size in pixels of the quiet zone around the QR code
	quietZone int

	moduleGap float64

	// TODO: Add better support of customization with single context (all colors, gaps, etc.) for both module and finders drawers
	moduleDrawer shapes.ModuleDrawer
	finderDrawer shapes.FinderDrawConfig

	// TODO
	// halftoneImg is the halftone image for the output image.
	//halftoneImg image.Image
}

var (
	color_WHITE = ParseFromHex("#ffffff")
	color_BLACK = ParseFromHex("#000000")
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
