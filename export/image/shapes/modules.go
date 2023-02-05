package shapes

import (
	"github.com/fogleman/gg"
)

type ModuleNeighbours struct {
	N, S, W, E bool
}

// ModuleDrawContext is a context for drawing single module on the image
type ModuleDrawContext struct {
	*gg.Context

	// X and Y is the coordinates of the top left corner of module
	X, Y float64
	// ModSize is the size of single module
	ModSize int
	// Gap should be subtracted from ModSize to keep padding between modules
	// Also, modules should be shifter by half of this value to properly centered
	//
	// This rule may be ignored if you need to connect modules in some specific way
	Gap float64

	// Neighbours contains a list of values for nearby neighbours
	Neighbours *ModuleNeighbours
}

// ModuleDrawer is a function that should Draw a single module
type ModuleDrawer struct {
	Draw            func(ctx *ModuleDrawContext)
	NeedsNeighbours bool
}

// SquareModuleShape draws simple square as module
func SquareModuleShape() ModuleDrawer {
	return RoundedModuleShape(0, false)
}

// RoundedModuleShape draws module as square with rounded corners.
//
// Supplied value is clamped between 0 (no roundness) and 0.5 (circle shape)
//
// connected determines whether modules will be connected between each other
//
// Note: when connected is true, this shape does not respect gaps (thus, gap will only be applied to finders)
func RoundedModuleShape(borderRadius float64, connected bool) ModuleDrawer {
	if borderRadius > 0.5 {
		borderRadius = 0.5
	}
	if borderRadius < 0 {
		borderRadius = 0
	}

	if connected {
		return ModuleDrawer{
			NeedsNeighbours: true,
			Draw: func(ctx *ModuleDrawContext) {
				size := float64(ctx.ModSize)
				n := ctx.Neighbours

				ctx.drawRoundedWithSquareCorners(
					ctx.X, ctx.Y, size, size, size*borderRadius,
					[4]bool{
						// Top right
						n.N || n.E,
						// Bottom right
						n.E || n.S,
						// Bottom left
						n.S || n.W,
						// Top left
						n.W || n.N,
					},
				)
			},
		}
	}

	return ModuleDrawer{
		NeedsNeighbours: false,
		Draw: func(ctx *ModuleDrawContext) {
			size := float64(ctx.ModSize) - ctx.Gap
			ctx.DrawRoundedRectangle(ctx.X+ctx.Gap/2, ctx.Y+ctx.Gap/2, size, size, size*borderRadius)
		},
	}
}

func LineModuleShape(borderRadius float64, vertical bool) ModuleDrawer {
	if borderRadius > 0.5 {
		borderRadius = 0.5
	}
	if borderRadius < 0 {
		borderRadius = 0
	}

	return ModuleDrawer{
		NeedsNeighbours: true,
		Draw: func(ctx *ModuleDrawContext) {
			size := float64(ctx.ModSize)
			// For applying gap on opposite direction (horizontal for vertical and vice versa)
			gapSize := size - ctx.Gap
			n := ctx.Neighbours

			if vertical {
				// Apply gap only on X axis (to pad rows)
				ctx.drawRoundedWithSquareCorners(
					ctx.X+ctx.Gap/2, ctx.Y, gapSize, size, gapSize*borderRadius,
					[4]bool{
						n.N, n.S, n.S, n.N,
					},
				)
				return
			}

			// Apply gap only on Y axis (to pad lines)
			ctx.drawRoundedWithSquareCorners(
				ctx.X, ctx.Y+ctx.Gap/2, size, gapSize, gapSize*borderRadius,
				[4]bool{
					n.E, n.E, n.W, n.W,
				},
			)
		},
	}
}

// Fixme: Can it be refactored?

// drawRoundedWithSquareCorners accept x and y position of rectangle, its width and height and the border radius for active corners
//
//	that are defined as true value in corners array, starting from top right corner and going clockwise
func (ctx *ModuleDrawContext) drawRoundedWithSquareCorners(x, y, w, h, r float64, squaredCorners [4]bool) {
	left, edgeLeft, edgeRight, right := x, x+r, x+w-r, x+w
	top, edgeTop, edgeBottom, bottom := y, y+r, y+h-r, y+h

	ctx.NewSubPath()
	ctx.MoveTo(edgeLeft, top)

	// Top & top right corner
	ctx.LineTo(edgeRight, top)
	if squaredCorners[0] {
		ctx.LineTo(right, top)
		ctx.LineTo(right, edgeTop)
	} else {
		ctx.drawArcDeg(edgeRight, edgeTop, r, 270, 360)
	}

	// Right & bottom right corner
	ctx.LineTo(right, edgeBottom)
	if squaredCorners[1] {
		ctx.LineTo(right, bottom)
		ctx.LineTo(edgeRight, bottom)
	} else {
		ctx.drawArcDeg(edgeRight, edgeBottom, r, 0, 90)
	}

	// Bottom & bottom right corner
	ctx.LineTo(edgeLeft, bottom)
	if squaredCorners[2] {
		ctx.LineTo(left, bottom)
		ctx.LineTo(left, edgeBottom)
	} else {
		ctx.drawArcDeg(edgeLeft, edgeBottom, r, 90, 180)
	}

	// Left & top left corner
	ctx.LineTo(left, edgeTop)
	if squaredCorners[3] {
		ctx.LineTo(left, top)
		ctx.LineTo(edgeLeft, top)
	} else {
		ctx.drawArcDeg(edgeLeft, edgeTop, r, 180, 270)
	}

	ctx.ClosePath()
}

func (ctx *ModuleDrawContext) drawArcDeg(x, y, r, angle1, angle2 float64) {
	ctx.DrawArc(x, y, r, gg.Radians(angle1), gg.Radians(angle2))
}
