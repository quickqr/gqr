package image

import (
	"github.com/fogleman/gg"
	"github.com/quickqr/gqr"
	"github.com/quickqr/gqr/export/image/imgkit"
	"github.com/quickqr/gqr/export/image/shapes"
	"golang.org/x/image/draw"
	"image"
	"image/color"
)

const logoSizeRatio float64 = 0.2

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
	//dc.SetColor(o.backgroundColor)

	dc.DrawRectangle(0, 0, float64(o.size), float64(o.size))
	dc.Fill()

	// Draw QR code data.
	qr := e.drawQR(&mat, actualSize)
	dc.DrawImage(qr, o.quietZone, o.quietZone)

	// Fixme: pieces of modules can be seen behind space container
	// TODO: find way to hide modules that are hidden by white space more than 80-90% (
	if o.logo != nil {
		// logo will automatically rescale to the size of QR code
		containerWidth := float64(actualSize) * logoSizeRatio
		imageWidth := int(containerWidth)

		if o.spaceAroundLogo {
			imageWidth = int(containerWidth * 0.8)

			//dc.SetColor(o.backgroundColor)
			center := (float64(o.size) - containerWidth) / 2
			dc.DrawRectangle(center, center, containerWidth, containerWidth)
			dc.Fill()
		}

		scaled := imgkit.Scale(o.logo, image.Rect(0, 0, imageWidth, imageWidth), nil)
		//should icon upper-left to start
		dc.DrawImage(scaled, (o.size-imageWidth)/2, (o.size-imageWidth)/2)
	}

	return dc.Image()

}

// TODO:
// - Add custom shapes for modules and finders, draw finders after other modules were drawn
// - Support for gradient (direction, set of colors)
// - Reimplement halftones capability from the original library

// drawQR draws pixel-perfect modules to avoid gaps between modules by ceiling width of module up and then
// scaling down image with data to actual QR width
func (e *Exporter) drawQR(mat *gqr.Matrix, requiredSize int) image.Image {
	// This line ceils division result without the need of converting everything to floats.
	// https://stackoverflow.com/a/2745086
	modSize := float64((e.options.size + mat.Width() - 1) / mat.Width())
	size := int(modSize) * mat.Width()
	// Apply gaps after real image size was calculated
	gap := modSize * e.options.moduleGap
	//whitespaceSize := float64(requiredSize) * logoSizeRatio
	//center := requiredSize / 2
	//emptyZone := (center - whitespaceSize)

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

	if e.options.gradientConfig != nil {
		c := e.options.gradientConfig
		size := float64(size)
		grad := dirToGradient(c.direction, size, size)
		l := len(c.colors)
		lastIdx := l - 1

		step := 1 / float64(l)

		// Set first color at 0
		grad.AddColorStop(0, c.colors[0])
		// Colors between 0 and last
		for i, v := range c.colors[1:lastIdx] {
			grad.AddColorStop(step*float64(i+1), v)
		}
		// Set last color at 1
		grad.AddColorStop(1, c.colors[lastIdx])

		dc.SetFillStyle(grad)
	} else {
		dc.SetColor(ctx.Color)
	}

	// iterate the matrix to Draw each pixel
	mat.Iterate(gqr.IterDirection_ROW, func(x int, y int, v gqr.QRValue) {
		// Finders are drawn separately
		if !v.IsSet() || v.Type() == gqr.QRType_FINDER {
			return
		}

		ctx.X = float64(x)*(ctx.ModSize) + gap/2
		ctx.Y = float64(y)*(ctx.ModSize) + gap/2

		if e.options.spaceAroundLogo {

		}

		e.options.drawModuleFn(ctx)

	})
	// Draw modules to screen
	dc.Fill()

	e.drawFinders(dc, modSize, gap)

	return imgkit.Scale(dc.Image(), image.Rect(0, 0, requiredSize, requiredSize), draw.ApproxBiLinear)
}

func (e *Exporter) drawFinders(dc *gg.Context, modSize float64, gap float64) {
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
