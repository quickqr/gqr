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
				//halfGap := ctx.Gap / 2

				if ctx.Neighbours != nil {
					ctx.drawConnectedRect(ctx.X, ctx.Y, size, size, size*borderRadius)
					return
				}

				size -= ctx.Gap
				ctx.DrawRoundedRectangle(ctx.X+ctx.Gap/2, ctx.Y+ctx.Gap/2, size, size, size*borderRadius)
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

// TODO: Find a way to refactor
func (ctx *ModuleDrawContext) drawConnectedRect(x, y, w, h, r float64) {
	n := ctx.Neighbours
	left, edgeLeft, edgeRight, right := x, x+r, x+w-r, x+w
	top, edgeTop, edgeBottom, bottom := y, y+r, y+h-r, y+h

	ctx.NewSubPath()
	ctx.MoveTo(edgeLeft, top)

	// Top & top right corner
	ctx.LineTo(edgeRight, top)
	if n.N || n.E {
		ctx.LineTo(right, top)
		ctx.LineTo(right, edgeTop)
	} else {
		ctx.drawArcDeg(edgeRight, edgeTop, r, 270, 360)
	}

	// Right & bottom right corner
	ctx.LineTo(right, edgeBottom)
	if n.S || n.E {
		ctx.LineTo(right, bottom)
		ctx.LineTo(edgeRight, bottom)
	} else {
		ctx.drawArcDeg(edgeRight, edgeBottom, r, 0, 90)
	}

	// Bottom & bottom right corner
	ctx.LineTo(edgeLeft, bottom)
	if n.S || n.W {
		ctx.LineTo(left, bottom)
		ctx.LineTo(left, edgeBottom)
	} else {
		ctx.drawArcDeg(edgeLeft, edgeBottom, r, 90, 180)
	}

	// Left & top left corner
	ctx.LineTo(left, edgeTop)
	if n.N || n.W {
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
