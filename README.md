# gqr
`gqr` is a QR code generation library focused on creating customizable QR codes.

This project is a fork of [github.com/yeqown/go-qrcode](http://github.com/yeqown/go-qrcode) with updated API (see Changes)
So big thanks to the initial project! 🙏

# Features
Features marked with `*` are inherited from `go-qrcode`
- [X] Normally generate QR code across version 1 to version 40 `*`
- [X] Automatically analyze QR version by source text `*`
- [X] Applying image size with `WithImageSize`
- [X] Applying padding for QR code with `WithQuietZone`
- [X] Full customization of shapes:  
  - [X] Customize modules with `WithModuleShape` 
    - Squares 
    - Rounded, 
    - Connected (no gaps between adjacent modules)
    - Lines
  - [X] Customize finders (eyes) with `WithFinderShape`
  - [X] "Connected" rounded modules with go gaps between
- [X] Setting colors `*`: `WithBgColor(Hex)`, `WithFgColor(Hex)`
  - [X] Added support for hex colors with alpha channel (8 digits)
  - [X] Gradient with customizable direction via `WithGradient(dir, colors...)`.  
    - Overrides foreground color
    - Colors will be automatically placed evenly on the gradient map
- [X] Advanced logo placing:
  - [X] `WithLogo` places image at center and rescales it to be 1/5 of qr code
  - [X] `WithSpaceAroundLogo` add white space at the center so logo is not placed on top of modules
- [ ] Support Halftones

# Install
```bash
go get -u github.com/quickqr/gqr
```

# Using
[Check other examples](./examples)
```go
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
```

# Showcase
All of these pictures are generated by programs in [examples](./examples):
<div float="left">
<img src="./assets/main.png" alt="main" width="300">
<img src="./assets/lines.png" alt="lines" width="300">
<img src="./assets/default.png" alt="default" width="300">
<img src="./assets/invert.png" alt="inverted" width="300">
<img src="assets/custom-shapes.png" alt="inverted" width="300">
</div>