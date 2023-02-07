package image

import (
	"github.com/fogleman/gg"
	"github.com/quickqr/gqr"
	"github.com/quickqr/gqr/export/image/imgkit"
	"github.com/quickqr/gqr/export/image/shapes"
	"golang.org/x/image/draw"
	"image"
	"math"
)

const logoSizeRatio float64 = 0.2

// conditionalRound performs a/b division and rounds up only if fraction part is more than bound
func conditionalRound(a float64, floor bool) float64 {
	if floor {
		return math.Floor(a)
	}
	return math.Ceil(a)
}

// Exporter exports gqr.Matrix to image.Image
type Exporter struct {
	options *exportOptions
}

// NewExporter creates new Exporter. (see DefaultImageOptions)
func NewExporter(opts ...ExportOption) Exporter {
	dst := DefaultImageOptions

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
	qr := e.drawQR(&mat, actualSize)
	dc.DrawImage(qr, o.quietZone, o.quietZone)

	if o.logo != nil {
		// rescale logo relative to size of the QR code
		containerWidth := int(float64(actualSize) * logoSizeRatio)

		imageWidth := int(containerWidth)
		if o.spaceAroundLogo {
			imageWidth = int(float64(containerWidth) * 0.8)
		}

		scaled := imgkit.Scale(o.logo, image.Rect(0, 0, imageWidth, imageWidth), nil)
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
	modSize := (e.options.size + mat.Width() - 1) / mat.Width()
	size := modSize * mat.Width()
	// Apply gaps after real image size was calculated
	gap := float64(modSize) * e.options.moduleGap
	e.clearEmptyZone(mat, size, modSize)

	dc := gg.NewContext(size, size)

	// qrcode block draw context
	ctx := &shapes.ModuleDrawContext{
		Context: dc,
		X:       0.0,
		Y:       0.0,
		ModSize: modSize,
		Gap:     gap,
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
		dc.SetColor(e.options.foregroundColor)
	}

	// iterate the matrix to Draw each pixel
	mat.Iterate(gqr.IterDirection_ROW, func(x int, y int, v gqr.QRValue) {
		// Finders are drawn separately
		if !v.IsSet() || v.Type() == gqr.QRType_FINDER {
			return
		}

		ctx.X = float64(x * ctx.ModSize)
		ctx.Y = float64(y * ctx.ModSize)

		if e.options.moduleDrawer.NeedsNeighbours {
			ctx.Neighbours = &shapes.ModuleNeighbours{
				N: mat.ValueAtClamped(x, y-1).IsSet(),
				S: mat.ValueAtClamped(x, y+1).IsSet(),
				W: mat.ValueAtClamped(x-1, y).IsSet(),
				E: mat.ValueAtClamped(x+1, y).IsSet(),
			}
		}

		e.options.moduleDrawer.Draw(ctx)
	})
	// Draw modules to screen
	dc.Fill()

	e.drawFinders(dc, float64(modSize), gap)

	return imgkit.Scale(dc.Image(), image.Rect(0, 0, requiredSize, requiredSize), draw.ApproxBiLinear)
}

// clearEmptyZone finds modules that intersect with logo container and hide them by unsetting
func (e *Exporter) clearEmptyZone(mat *gqr.Matrix, imageSize, modSize int) {
	if !e.options.spaceAroundLogo {
		return
	}

	emptyZone := e.getEmptyZone(imageSize)
	min := emptyZone.Min.X / modSize
	max := (emptyZone.Max.X + modSize - 1) / modSize

	for x := min; x <= max; x++ {
		for y := min; y <= max; y++ {
			xRect := x * modSize
			yRect := y * modSize

			if image.Rect(xRect, yRect, xRect+modSize, yRect+modSize).Overlaps(emptyZone) {
				_ = mat.Set(x, y, gqr.QRValue_DATA_V0)
			}

		}
	}
}

func (e *Exporter) getEmptyZone(imageSize int) image.Rectangle {
	if e.options.logo != nil && e.options.spaceAroundLogo {
		center := imageSize / 2
		// Get half of the empty zone size
		halfSize := int(float64(imageSize) * logoSizeRatio / 2)
		start := center - halfSize
		end := center + halfSize

		return image.Rect(start, start, end, end)
	}

	return image.Rect(0, 0, 0, 0)
}

func (e *Exporter) drawFinders(dc *gg.Context, modSize float64, gap float64) {
	finderSize := modSize * gqr.FINDER_SIZE
	modSize -= gap
	// This will not draw second WhiteSpace and cut a "hole" of specified shape inside finderDrawer.Outer
	dc.SetFillRuleEvenOdd()

	// Placing outer shapes
	placeFinderShapes(dc, e.options.finderDrawer.Outer, finderSize, modSize, 0)

	// Placing space between outer and inner shapes
	placeFinderShapes(dc, e.options.finderDrawer.WhiteSpace, finderSize-modSize*2, modSize, modSize)

	innerSize := finderSize / 2
	placeFinderShapes(dc, e.options.finderDrawer.Inner, innerSize, modSize, (finderSize-innerSize)/2)

	dc.Fill()
}

func placeFinderShapes(ctx *gg.Context, f shapes.FinderShapeDrawer, size float64, modSize float64, offset float64) {
	finderSize := size

	// Top Left
	f(ctx, offset, offset, finderSize, modSize)
	// Top right
	f(ctx, -offset+float64(ctx.Width())-size, offset, finderSize, modSize)
	// Bottom left
	f(ctx, offset, -offset+float64(ctx.Width())-size, finderSize, modSize)
}
