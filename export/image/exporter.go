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
	qrData := e.getQRImage(&mat, actualSize)
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

// getQRImage draws pixel-perfect modules to avoid gaps between modules by ceiling width of module up and then
// scaling down image with data to actual QR width
func (e *Exporter) getQRImage(mat *gqr.Matrix, requiredSize int) image.Image {
	// This line ceils division result without the need of converting everything to floats.
	// https://stackoverflow.com/a/2745086
	modSize := float64((e.options.size + mat.Width() - 1) / mat.Width())

	size := int(modSize) * mat.Width()

	// Apply gaps after real image size was calculated
	gap := modSize * e.options.moduleGap

	dc := gg.NewContext(size, size)

	// qrcode block draw context
	ctx := &shapes.DrawContext{
		Context: dc,
		X:       0.0,
		Y:       0.0,
		ModSize: modSize,
		Gap:     gap,
		Color:   e.options.foregroundColor,
	}

	ctx.SetColor(ctx.Color)
	// iterate the matrix to Draw each pixel
	mat.Iterate(gqr.IterDirection_ROW, func(x int, y int, v gqr.QRValue) {
		// Finders are drawn separately
		if !v.IsSet() || v.Type() == gqr.QRType_FINDER {
			return
		}

		ctx.X = float64(x)*(ctx.ModSize) + gap/2
		ctx.Y = float64(y)*(ctx.ModSize) + gap/2

		e.options.drawModuleFn(ctx)

	})
	//// Draw modules to screen
	ctx.Fill()

	e.drawFinders(dc, modSize, gap)

	return imgkit.Scale(dc.Image(), image.Rect(0, 0, requiredSize, requiredSize), nil)
}

func (e *Exporter) drawFinders(dc *gg.Context, modSize float64, gap float64) {
	dc.SetColor(e.options.foregroundColor)
	finderSize := modSize * gqr.FINDER_SIZE
	modSize -= gap

	// Creating mask to cut out inside of outer shapes
	mask := gg.NewContext(dc.Width(), dc.Height())
	mask.SetColor(color.Black)
	placeFinderShapes(mask, e.options.drawFinder.WhiteSpace, finderSize, modSize)
	mask.Fill()
	_ = dc.SetMask(mask.AsMask())
	dc.InvertMask()

	// Placing outer shapes
	placeFinderShapes(dc, e.options.drawFinder.Outer, finderSize, 0)

	// Resetting mask to set inner shapes
	mask.Clear()
	_ = dc.SetMask(mask.AsMask())

	placeFinderShapes(dc, e.options.drawFinder.WhiteSpace, finderSize, 2*modSize)
}

func placeFinderShapes(ctx *gg.Context, f shapes.FinderShapeDrawer, size float64, offset float64) {
	offsetSize := size - offset*2
	f(ctx, offset, offset, offsetSize)
	f(ctx, offset, offset+float64(ctx.Width())-size, offsetSize)
	f(ctx, offset+float64(ctx.Width())-size, offset, offsetSize)
}
