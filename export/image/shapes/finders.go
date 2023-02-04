package shapes

import (
	"github.com/fogleman/gg"
)

// FinderShapeDrawer is general function for drawing finder shapes
//
// This function should Draw any shape with top left corned on cords x and y with supplied size.
// FinderShapeDrawer should only Draw a path, but not fill anything, as it's done when calling this function
type FinderShapeDrawer = func(ctx *gg.Context, x float64, y float64, size float64, modSize float64)

// FinderDrawConfig contains 3 FinderShapeDrawer functions, this allows to customize look of the inner and outer borders.
type FinderDrawConfig struct {
	// Outer shape places outer container
	Outer FinderShapeDrawer
	// WhiteSpace used to create a mask to cut out inner part of the Outer shape
	WhiteSpace FinderShapeDrawer
	// Inner shape is drawn lastly, at the center of white space created by WhiteSpace shape
	Inner FinderShapeDrawer
}

// SquareFinderShape is default shape for finders. Draws simple squares
func SquareFinderShape() FinderDrawConfig {
	return RoundedFinderShape(0)
}

// RoundedFinderShape draws finders with rounded corners.
//
// borderRadius should be between 0 (square) and 0.5 (full circle). Other values are clamped to fit these requirements
func RoundedFinderShape(borderRadius float64) FinderDrawConfig {
	if borderRadius > 0.5 {
		borderRadius = 0.5
	}
	if borderRadius < 0 {
		borderRadius = 0
	}

	draw := func(ctx *gg.Context, x float64, y float64, size float64, modSize float64) {
		ctx.DrawRoundedRectangle(x, y, size, size, size*borderRadius)
	}

	return FinderDrawConfig{draw, draw, draw}
}
