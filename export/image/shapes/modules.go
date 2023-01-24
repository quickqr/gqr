package shapes

// TODO: Margin should be set inside of these functions because some of them need to control when to add margin
// (for example, when using "connectModules" in RoundedModuleShape

// SquareModuleShape draws simple rectangle as module
func SquareModuleShape() ModuleShapeDrawer {
	return func(ctx *DrawContext) {
		ctx.DrawRectangle(ctx.X, ctx.Y, ctx.Width, ctx.Width)
		ctx.SetColor(ctx.Color)
		ctx.Fill()
	}
}

// RoundedModuleShape draws module with rounded corners.
// Supplied value is clamped between 0 (no roundness) and 0.5 (circle shape)
// TODO: Add "connectModules" option to specify whether to connect modules when neighbouring
func RoundedModuleShape(borderRadius float64) ModuleShapeDrawer {
	if borderRadius > 0.5 {
		borderRadius = 0.5
	}
	if borderRadius < 0 {
		borderRadius = 0
	}

	return func(ctx *DrawContext) {
		ctx.DrawRoundedRectangle(ctx.X, ctx.Y, ctx.Width, ctx.Width, ctx.Width*borderRadius)
		ctx.SetColor(ctx.Color)
		ctx.Fill()
	}
}
