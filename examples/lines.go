package main

import (
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

	// Export QR code to image
	img := export.
		NewExporter(
			export.WithBgColorHex("#1f1f1f"),

			export.WithFinderShape(shapes.RoundedFinderShape(0.3)),
			export.WithModuleShape(shapes.LineModuleShape(0.5, false)),
			export.WithModuleGap(0.3),

			export.WithGradient(export.GradientDirectionLTR,
				export.ParseFromHex("#cc33ff"),
				export.ParseFromHex("#ff9900"),
				// You also can use any color.Color instance
			),
		).
		Export(*qr)

	// Save the image
	f, _ := os.OpenFile("../assets/lines.png", os.O_CREATE|os.O_WRONLY, 0644)
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
