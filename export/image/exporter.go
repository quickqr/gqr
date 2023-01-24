package image

import (
	"github.com/fogleman/gg"
	"github.com/quickqr/gqr"
	"github.com/quickqr/gqr/export/image/imgkit"
	"github.com/quickqr/gqr/export/image/shapes"
	"image"
	"image/color"
)

// Exporter exports gqr.Matrix to image.Image
type Exporter struct {
	options *imageOptions
}

// NewExporter creates new Exporter. (see defaultImageOptions)
func NewExporter(opts ...ImageOption) Exporter {
	dst := defaultImageOptions

	for _, opt := range opts {
		opt.apply(&dst)
	}

	return Exporter{
		options: &dst,
	}
}

// Export QR to Image.image
func (e Exporter) Export(mat gqr.Matrix) image.Image {
	o := e.options
	actualSize := o.size - o.quietZone*2
	dc := gg.NewContext(o.size, o.size)

	// Draw background
	dc.SetColor(o.backgroundColor)
	dc.DrawRectangle(0, 0, float64(o.size), float64(o.size))
	dc.Fill()

	// Draw QR code data.
	qrData := e.getDataImage(&mat, actualSize)
	dc.DrawImage(qrData, o.quietZone, o.quietZone)

	// TODO: Add support for logo image background container
	if o.logo != nil {
		// logo will automatically rescale to the size of QR code
		logoWidth := actualSize
		scaled := imgkit.Scale(o.logo, image.Rect(0, 0, logoWidth, logoWidth), nil)

		//should icon upper-left to start
		dc.DrawImage(scaled, (o.size-logoWidth)/2, (o.size-logoWidth)/2)
	}

	return dc.Image()

}

// TODO:
// - Add custom shapes for modules and finders, draw finders after other modules were drawn
// - Support for gradient (direction, set of colors)
// - Reimplement halftones capability from the original library

// getDataImage draws pixel-perfect modules to avoid gaps between modules by ceiling width of module up and then
// scaling down image with data to actual QR width
func (e *Exporter) getDataImage(mat *gqr.Matrix, actualSize int) image.Image {
	// This line ceils division result without the need of converting everything to floats.
	// https://stackoverflow.com/a/2745086
	modW := (e.options.size + mat.Width() - 1) / mat.Width()
	size := modW * mat.Width()
	dc := gg.NewContext(size, size)

	// qrcode block draw context
	ctx := &shapes.DrawContext{
		Context: dc,
		X:       0.0,
		Y:       0.0,
		Width:   modW,
		Height:  modW,
		Color:   color.Black,
	}

	// iterate the matrix to Draw each pixel
	mat.Iterate(gqr.IterDirection_ROW, func(x int, y int, v gqr.QRValue) {
		if v.Type() == gqr.QRType_FINDER {
			return
		}

		ctx.X = float64(x * ctx.Width)
		ctx.Y = float64(y * ctx.Width)
		//ctx.Width = modW - 2

		ctx.Color = e.options.qrValueToRGBA(v)

		e.options.drawModuleFn(ctx)

		// Fixme: add generic shapes back
		//ctx.DrawRectangle(ctx.X, ctx.Y, float64(ctx.Width), float64(ctx.Width))
		//ctx.SetColor(ctx.Color)
		//ctx.Fill()

		// FIXME: Should ignore Finders
		//switch typ := v.Type(); typ {
		//case gqr.QRType_FINDER:
		//	break
		//default:
		//}
	})

	return imgkit.Scale(dc.Image(), image.Rect(0, 0, actualSize, actualSize), nil)
}
