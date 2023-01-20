package image

import (
	"github.com/stretchr/testify/assert"
	"image/color"
	"testing"
)

func Test_BgColor_FgColor(t *testing.T) {
	oo := &defaultImageOptions

	// check
	assert.Equal(t, color_WHITE, oo.backgroundColor)
	assert.Equal(t, color_BLACK, oo.foregroundColor)

	// apply color
	WithBgColor(color_BLACK).apply(oo)
	assert.Equal(t, color_BLACK, oo.backgroundColor)
	assert.Equal(t, color_BLACK, oo.foregroundColor)

	// apply color
	WithBgColor(color_WHITE).apply(oo)
	WithFgColor(color_WHITE).apply(oo)
	assert.Equal(t, color_WHITE, oo.backgroundColor)
	assert.Equal(t, color_WHITE, oo.foregroundColor)

	WithFgColor(color_BLACK).apply(oo)
	assert.Equal(t, color_WHITE, oo.backgroundColor)
	assert.Equal(t, color_BLACK, oo.foregroundColor)
}

func Test_defaultOutputOption(t *testing.T) {
	oo := &defaultImageOptions

	// Apply
	rgba := color.RGBA{
		R: 123,
		G: 123,
		B: 123,
		A: 123,
	}
	WithBgColor(rgba).apply(oo)
	// assert
	assert.Equal(t, rgba, oo.backgroundColor)

	// check default
	oo2 := defaultImageOptions
	assert.NotEqual(t, oo2.backgroundColor, oo.backgroundColor)
}

func Test_WithQuietZoneSize(t *testing.T) {
	oo := &defaultImageOptions

	// zero parameter
	WithQuietZone(50).apply(oo)
	assert.Equal(t, 500, oo.quietZone)
}
