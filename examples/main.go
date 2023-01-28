package main

import (
	"github.com/quickqr/gqr"
	export "github.com/quickqr/gqr/export/image"
	"github.com/quickqr/gqr/export/image/shapes"
	"image"
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

	logoFile, _ := os.Open("./gopher.png")
	logo, _, _ := image.Decode(logoFile)

	// Export QR code to image
	img := export.
		NewExporter(
			export.WithBgColorHex("#1f1f1f"),

			export.WithLogo(logo),
			export.WithSpaceAroundLogo(),

			export.WithFinderShape(shapes.RoundedFinderShape(0.5)),
			export.WithModuleShape(shapes.RoundedModuleShape(0.5)),

			// Apply gap between modules
			export.WithModuleGap(0.1),
			// Size of the outputted image in pixels
			export.WithImageSize(512),
			// Padding around QR code
			// Note: actual QR code size will be (image size - quiet zone * 2)
			export.WithQuietZone(30),

			// Gradient for foreground with direction from Top Right to Bottom Left
			export.WithGradient(export.GradientDirectionTRBL,
				export.ParseFromHex("#cc33ff"),
				export.ParseFromHex("#ff9900"),
				// You also can use any color.Color instance
			),
		).
		Export(*qr)

	// Save the image
	f, _ := os.OpenFile("../assets/main.png", os.O_CREATE|os.O_WRONLY, 0644)
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
