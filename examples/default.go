package main

import (
	"github.com/quickqr/gqr"
	export "github.com/quickqr/gqr/export/image"
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
	img := export.NewExporter().Export(*qr)

	// Save the image
	f, _ := os.OpenFile("../assets/default.png", os.O_CREATE|os.O_WRONLY, 0644)
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
