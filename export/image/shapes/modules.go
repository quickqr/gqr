package shapes

func SquareModuleShape() ModuleShapeDrawer {
	return func(ctx *DrawContext) {
		ctx.DrawRectangle(ctx.X, ctx.Y, float64(ctx.Width), float64(ctx.Width))
		ctx.SetColor(ctx.Color)
		ctx.Fill()
	}
}

// RoundedModuleShape draws module with rounded corners. Border Radius is percent value of roundness between 0 and 50
// where 0 is no border radius and 50 is complete circle
func RoundedModuleShape(borderRadius uint8) ModuleShapeDrawer {
	if borderRadius > 50 {
		borderRadius = 50
	}

	rad := float64(borderRadius) / 100.0

	return func(ctx *DrawContext) {
		ctx.DrawRoundedRectangle(ctx.X, ctx.Y, float64(ctx.Width), float64(ctx.Width), float64(ctx.Width)*rad)
		ctx.SetColor(ctx.Color)
		ctx.Fill()
	}
}
