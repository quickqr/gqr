package shapes

import (
	"github.com/fogleman/gg"
)

// ModuleDrawContext is a context for drawing single module on the image
type ModuleDrawContext struct {
	*gg.Context

	// X and Y is the coordinates of the top left corner of module
	X, Y float64
	// ModSize is the size of single module
	ModSize int
	// Gap should be subtracted from ModSize to keep padding between modules
	//
	// This rule may be ignored if you need to connect modules in some specific way
	Gap float64
}

// ModuleDrawer is a function that should draw a single module
type ModuleDrawer = func(ctx *ModuleDrawContext)

// SquareModuleShape draws simple square as module
func SquareModuleShape() ModuleDrawer {
	return RoundedModuleShape(0)
}

// RoundedModuleShape draws module as square with rounded corners.
//
// Supplied value is clamped between 0 (no roundness) and 0.5 (circle shape)
func RoundedModuleShape(borderRadius float64) ModuleDrawer {
	if borderRadius > 0.5 {
		borderRadius = 0.5
	}
	if borderRadius < 0 {
		borderRadius = 0
	}

	// TODO: Implement "connectNeighbours" feature
	return func(ctx *ModuleDrawContext) {
		size := float64(ctx.ModSize) - ctx.Gap
		ctx.DrawRoundedRectangle(ctx.X, ctx.Y, size, size, size*borderRadius)
	}
}
