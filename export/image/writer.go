package image

import (
	"github.com/fogleman/gg"
	"github.com/quickqr/gqr"
	"github.com/quickqr/gqr/export/image/imgkit"
	"image"
	"image/color"
)

// Exporter exports gqr.Matrix to image.Image
type Exporter struct {
	options *outputImageOptions
}

// New creates new Exporter. (see defaultImageOptions)
func New(opts ...ImageOption) Exporter {
	dst := defaultImageOptions

	for _, opt := range opts {
		opt.apply(&dst)
	}

	return Exporter{
		options: &dst,
	}
}

// TODO:
// - Draw finders separately

// Export QR
func (e Exporter) Export(mat gqr.Matrix) image.Image {
	o := e.options

	dc := gg.NewContext(o.Size, o.Size)

	// draw background
	dc.SetColor(o.backgroundColor)
	dc.DrawRectangle(0, 0, float64(o.Size), float64(o.Size))
	dc.Fill()

	actualSize := o.Size - o.quietZone*2
	modWidth := float64(actualSize) / float64(mat.Width())

	// Fixme: should be avaliable as boolean like "touchesNeighbours" in shape interface
	// 1 pixel of padding to cover gaps that appear after floating point arithmetics
	modPad := 1.0

	// qrcode block draw context
	ctx := &DrawContext{
		Context: dc,
		x:       0.0,
		y:       0.0,
		w:       modWidth + modPad,
		h:       modWidth + modPad,
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
		ctx.x = float64(x)*modWidth + float64(o.quietZone)
		ctx.y = float64(y)*modWidth + float64(o.quietZone)
		ctx.color = o.qrValueToRGBA(v)

		switch typ := v.Type(); typ {
		case gqr.QRType_FINDER:
			// Pass finders since they're drawn separately
			break
		default:
			// Fixme: add generic shapes back
			ctx.DrawRectangle(ctx.x, ctx.y, float64(ctx.w), float64(ctx.w))
			ctx.SetColor(ctx.color)
			ctx.Fill()
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