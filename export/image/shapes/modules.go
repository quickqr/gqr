package shapes

import (
	"github.com/fogleman/gg"
	"image/color"
)

// DrawContext is a context for drawing single module on the image
type DrawContext struct {
	*gg.Context

	X, Y    float64
	ModSize float64
	Gap     float64

	Color color.Color
}

type ModuleDrawer = func(ctx *DrawContext)

// SquareModuleShape draws simple square as module
func SquareModuleShape() ModuleDrawer {
	return RoundedModuleShape(0)
}

// RoundedModuleShape draws module as square with rounded corners.
// Supplied value is clamped between 0 (no roundness) and 0.5 (circle shape)
func RoundedModuleShape(borderRadius float64) ModuleDrawer {
	if borderRadius > 0.5 {
		borderRadius = 0.5
	}
	if borderRadius < 0 {
		borderRadius = 0
	}

	// TODO: Implement "connectNeighbours" feature
	return func(ctx *DrawContext) {
		size := ctx.ModSize - ctx.Gap
		ctx.DrawRoundedRectangle(ctx.X, ctx.Y, size, size, size*borderRadius)
	}
}
