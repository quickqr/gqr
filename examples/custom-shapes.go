package main

import (
	"github.com/fogleman/gg"
	"github.com/quickqr/gqr"
	export "github.com/quickqr/gqr/export/image"
	"github.com/quickqr/gqr/export/image/shapes"
	"image/png"
	"log"
	"os"
)

func main() {
	qr, e := gqr.NewWith(
		"https://github.com/quickqr/gqr",
	)

	if e != nil {
		log.Fatal(e)
	}

	finderDrawer := shapes.SquareFinderShape()
	finderDrawer.Inner = innerPolygon
	moduleDrawer := drawModule

	// Export QR code to image
	img := export.NewExporter(
		export.WithBgColorHex("#1f1f1f"),
		export.WithGradient(export.GradientDirectionLTR,
			export.ParseFromHex("#00d4ff"),
			export.ParseFromHex("#3037ad"),
		),
		export.WithFinderShape(finderDrawer),
		export.WithModuleShape(moduleDrawer),
	).
		Export(*qr)

	// Save the image
	f, _ := os.OpenFile("../assets/custom-shapes.png", os.O_CREATE|os.O_WRONLY, 0644)
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

// You can refer to https://github.com/quickqr/gqr/tree/main/export/image/shapes to see more inspiration
func drawModule(ctx *shapes.ModuleDrawContext) {
	size := float64(ctx.ModSize) - ctx.Gap
	rad := size / 2
	ctx.DrawRegularPolygon(6, ctx.X+rad, ctx.Y+rad, rad, 90)
}

func innerPolygon(dc *gg.Context, x, y, size, modSize float64) {
	dc.DrawRegularPolygon(6, x+size/2, y+size/2, size/2, 0)
}
