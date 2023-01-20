package image

import (
	"github.com/fogleman/gg"
	"github.com/quickqr/gqr"
	"github.com/quickqr/gqr/writer/image/imgkit"
	"image"
	"image/color"
)

// ImageExporter exports gqr.Matrix to image.Image
type ImageExporter struct {
	options *ImageOptions
}

// New creates new ImageExporter with default options. (see DEFAULT_IMAGE_OPTIONS)
func New() ImageExporter {
	return NewWithOptions(DEFAULT_IMAGE_OPTIONS)
}

// NewWithOptions creates ImageExporter with custom options. (see ImageOptions)
func NewWithOptions(opt ImageOptions) ImageExporter {
	return ImageExporter{
		options: &opt,
	}
}

// TODO:
// - Draw finders separately

// Export QR
func (e ImageExporter) Export(mat gqr.Matrix) image.Image {
	o := e.options

	dc := gg.NewContext(o.Size, o.Size)

	// draw background
	dc.SetColor(o.BackgroundColor)
	dc.DrawRectangle(0, 0, float64(o.Size), float64(o.Size))
	dc.Fill()

	actualSize := o.Size - o.BorderWidth*2
	modWidth := float64(actualSize) / float64(mat.Width())
	// qrcode block draw context
	ctx := &DrawContext{
		Context: dc,
		x:       0.0,
		y:       0.0,
		w:       modWidth,
		h:       modWidth,
		color:   color.Black,
	}

	//var (
	//	halftoneImg image.Image
	//	halftoneW   = float64(opt.qrBlockWidth()) / 3.0
	//)
	//if opt.halftoneImg != nil {
	//	halftoneImg = imgkit.Binaryzation(
	//		imgkit.Scale(opt.halftoneImg, image.Rect(0, 0, mat.Width()*3, mat.Width()*3), nil),
	//		60,
	//	)
	//
	//	//_ = imgkit.Save(halftoneImg, "mask.jpeg")
	//}

	// iterate the matrix to Draw each pixel
	mat.Iterate(gqr.IterDirection_ROW, func(x int, y int, v gqr.QRValue) {
		ctx.x = modWidth*float64(x) + float64(o.BorderWidth)
		ctx.y = modWidth*float64(y) + float64(o.BorderWidth)
		ctx.color = o.qrValueToRGBA(v)

		switch typ := v.Type(); typ {
		case gqr.QRType_FINDER:
			// Pass finders since they're drawn separately
			break
		case gqr.QRType_DATA:
			// TODO: Fix halftones
			//if halftoneImg == nil {
			//	shape.Draw(ctx)
			//	return
			//}
			//
			//ctx2 := &DrawContext{
			//	Context: ctx.Context,
			//	w:       int(halftoneW),
			//	h:       int(halftoneW),
			//}
			//// only halftone image enabled and current block is Data.
			//for i := 0; i < 3; i++ {
			//	for j := 0; j < 3; j++ {
			//		ctx2.x, ctx2.y = ctx.x+float64(i)*halftoneW, ctx.y+float64(j)*halftoneW
			//		ctx2.color = halftoneImg.At(x*3+i, y*3+j)
			//		if i == 1 && j == 1 {
			//			ctx2.color = ctx.color
			//			// only center block keep the origin color.
			//		}
			//		shape.Draw(ctx2)
			//	}
			//}
		default:
			// Fixme: add generic shapes back
			ctx.DrawRectangle(ctx.x, ctx.y, float64(ctx.w), float64(ctx.h))
			ctx.SetColor(ctx.color)
			ctx.Fill()
		}
	})

	if o.Logo != nil {
		// Logo will automatically rescale to the size of QR code
		logoWidth := actualSize
		scaled := imgkit.Scale(o.Logo, image.Rect(0, 0, logoWidth, logoWidth), nil)

		//should icon upper-left to start
		dc.DrawImage(scaled, (o.Size-logoWidth)/2, (o.Size-logoWidth)/2)
	}

	return dc.Image()

}

// Attribute contains basic information of generated image.
type Attribute struct {
	// width and height of image
	W, H   int
	Border int
}
