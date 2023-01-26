package shapes

import (
	"github.com/fogleman/gg"
)

type FinderShapeDrawer = func(ctx *gg.Context, x float64, y float64, size float64)
type FinderDrawConfig struct {
	Outer      FinderShapeDrawer
	WhiteSpace FinderShapeDrawer
	Inner      FinderShapeDrawer
}

func SquareFinderShape() FinderDrawConfig {
	return RoundedFinderShape(0)
}

func RoundedFinderShape(borderRadius float64) FinderDrawConfig {
	if borderRadius > 0.5 {
		borderRadius = 0.5
	}
	if borderRadius < 0 {
		borderRadius = 0
	}

	draw := func(ctx *gg.Context, x float64, y float64, size float64) {
		ctx.DrawRoundedRectangle(x, y, size, size, size*borderRadius)
		ctx.Fill()
	}

	return FinderDrawConfig{draw, draw, draw}
}
