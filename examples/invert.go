package main

import (
	"github.com/quickqr/gqr"
	export "github.com/quickqr/gqr/export/image"
	"github.com/quickqr/gqr/export/image/shapes"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	qr, e := gqr.NewWith(
		"Hello!",
	)

	if e != nil {
		log.Fatal(e)
	}

	// Export QR code to image
	img := export.
		NewExporter(
			export.WithBgColor(color.Black),
			export.WithFgColor(color.White),
			export.WithQuietZone(100),
			// Works better with rounded module shapes
			export.WithModuleGap(0.1),
			export.WithModuleShape(shapes.RoundedModuleShape(0.2)),
		).
		Export(*qr)

	// Save the image
	f, _ := os.OpenFile("../assets/invert.png", os.O_CREATE|os.O_WRONLY, 0644)
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
