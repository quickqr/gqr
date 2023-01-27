package image

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ApplyOptions(t *testing.T) {
	desired := imageOptions{
		backgroundColor: color_BLACK,               // white
		foregroundColor: parseFromHex("#1f1f1f00"), // black
		size:            1024,
		quietZone:       50,
		moduleGap:       0.2,
	}

	e := NewExporter(
		WithBgColor(color_BLACK),
		WithFgColorHex("#1f1f1f00"),
		WithImageSize(1024),
		WithQuietZone(50),
		WithModuleGap(0.2),
	)

	assert.Equal(t, desired.backgroundColor, e.options.backgroundColor)
	assert.Equal(t, desired.foregroundColor, e.options.foregroundColor)
	assert.Equal(t, desired.size, e.options.size)
	assert.Equal(t, desired.quietZone, e.options.quietZone)
	assert.Equal(t, desired.moduleGap, e.options.moduleGap)
}
