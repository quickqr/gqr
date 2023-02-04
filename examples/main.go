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
		// These are defaults, can be omitted (as it is in other examples)
		gqr.WithErrorCorrectionLevel(gqr.ErrorCorrectionHighest),
		gqr.WithEncodingMode(gqr.EncModeAuto),
		// Uncomment to force QR code version 23.
		// Note: you might get an error if supplied length is more than capacity of a specified version
		// gqr.WithVersion(23)
	)

	if e != nil {
		log.Fatal(e)
	}

	logoFile, _ := os.Open("./gopher.png")
	logo, _, _ := image.Decode(logoFile)

	// Export QR code to image
	img := export.
		NewExporter(
			export.WithBgColorHex("#ffffff"),

			export.WithLogo(logo),
			export.WithSpaceAroundLogo(),

			export.WithFinderShape(shapes.RoundedFinderShape(0.3)),
			export.WithModuleShape(shapes.RoundedModuleShape(0.5, true)),

			// Apply gap between modules (note: this particular example will not use it because of connected modules)
			// You can see this working as expected in invert.go and  custom-shapes.go
			export.WithModuleGap(0.1),
			// Size of the outputted image in pixels
			export.WithImageSize(1024),
			// Padding around QR code
			// Note: actual QR code size will be (image size - quiet zone * 2)
			export.WithQuietZone(60),

			// Gradient for foreground with direction from Top Right to Bottom Left
			export.WithGradient(export.GradientDirectionLTR,
				export.ParseFromHex("#336FE1"),
				export.ParseFromHex("#2799C9"),
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
