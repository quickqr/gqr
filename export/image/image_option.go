package image

import (
	"fmt"
	"github.com/quickqr/gqr"
	"image"
	"image/color"
)

var DEFAULT_IMAGE_OPTIONS = ImageOptions{
	BackgroundColor: color_WHITE, // white
	ForegroundColor: color_BLACK, // black
	Logo:            nil,
	Size:            512,
	BorderWidth:     30,
}

// ImageOptions to output QR code image
type ImageOptions struct {
	// BackgroundColor is the background color of the QR code image.
	BackgroundColor color.RGBA

	// ForegroundColor is the foreground color of the QR code.
	ForegroundColor color.RGBA

	// TODO: Add more customization for logo, probably move into own LogoOptions

	// Logo this icon image would be put the center of QR Code image
	Logo image.Image

	// Size in pixel of output image
	// Note: Actual size of the QR code will be equal to Size - BorderWidth
	Size int

	// BorderWidth is the size in pixels of the quiet zone around the QR code
	BorderWidth int

	// TODO
	// halftoneImg is the halftone image for the output image.
	//halftoneImg image.Image
}

func (oo *ImageOptions) qrValueToRGBA(v gqr.QRValue) (rgba color.RGBA) {
	if v.IsSet() {
		rgba = oo.ForegroundColor
		return rgba
	}

	return rgba
}

// preCalculateAttribute this function must reference to draw function.
func (oo *ImageOptions) preCalculateAttribute(dimension int) *Attribute {
	if oo == nil {
		return nil
	}

	//top, right, bottom, left := oo.borderWidths[0], oo.borderWidths[1], oo.borderWidths[2], oo.borderWidths[3]
	return &Attribute{
		W:      oo.Size + oo.BorderWidth*2,
		H:      oo.Size + oo.BorderWidth*2,
		Border: oo.BorderWidth,
	}
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
