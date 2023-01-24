package shapes

import (
	"github.com/fogleman/gg"
	"image/color"
)

// DrawContext is a rectangle area
type DrawContext struct {
	*gg.Context

	X, Y          float64
	Width, Height int

	Color color.Color
}

type ModuleShapeDrawer = func(ctx *DrawContext)
