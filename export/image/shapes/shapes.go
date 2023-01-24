package shapes

import (
	"github.com/fogleman/gg"
	"image/color"
)

// DrawContext is a rectangle area
type DrawContext struct {
	*gg.Context

	X, Y          float64
	Width, Height float64

	Color color.Color
}

// TODO: Access to matrix, gap info
type ModuleShapeDrawer = func(ctx *DrawContext)
