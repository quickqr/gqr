package image

import (
	"github.com/stretchr/testify/assert"
	"image/color"
	"testing"
)

func Test_BgColor_FgColor(t *testing.T) {
	oo := &DefaultImageOptions

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

func Test_WithImageSize(t *testing.T) {
	oo := DefaultImageOptions

	WithImageSize(20).apply(&oo)
	// assert
	assert.Equal(t, 20, oo.size)

	WithImageSize(-10).apply(&oo)
	// assert ignore invalid
	assert.Equal(t, 20, oo.size)
}

func Test_WithModuleGap(t *testing.T) {
	oo := DefaultImageOptions

	WithModuleGap(0.5).apply(&oo)
	// assert
	assert.Equal(t, 0.5, oo.moduleGap)

	WithModuleGap(-1).apply(&oo)
	WithModuleGap(2).apply(&oo)
	// assert ignore invalid
	assert.Equal(t, 0.5, oo.moduleGap)
}

func Test_defaultOutputOption(t *testing.T) {
	oo := DefaultImageOptions

	// Apply
	rgba := color.RGBA{
		R: 123,
		G: 123,
		B: 123,
		A: 123,
	}
	WithBgColor(rgba).apply(&oo)
	// assert
	assert.Equal(t, rgba, oo.backgroundColor)

	// check default
	assert.NotEqual(t, DefaultImageOptions.backgroundColor, oo.backgroundColor)
}

func Test_WithQuietZoneSize(t *testing.T) {
	oo := DefaultImageOptions

	// zero parameter
	WithQuietZone(50).apply(&oo)
	assert.Equal(t, 50, oo.quietZone)
}

func Test_WithLogoScale(t *testing.T) {
	oo := DefaultImageOptions

	// Too small value
	WithLogoScale(0.35).apply(&oo)
	assert.Equal(t, DefaultImageOptions.logoScale, oo.logoScale)
	// Too big value
	WithLogoScale(1.25).apply(&oo)
	assert.Equal(t, DefaultImageOptions.logoScale, oo.logoScale)

	// Ok value
	WithLogoScale(0.65).apply(&oo)
	assert.Equal(t, 0.65, oo.logoScale)
}
