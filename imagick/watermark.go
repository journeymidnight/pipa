package imagick

import (
	"gopkg.in/gographics/imagick.v3/imagick"
)

type Watermark struct {
	XMargin      int
	YMargin      int
	Transparency float64
	Picture      *imagick.MagickWand //watermark is picture
}

type Text struct {
	text     string
	color    string
	front    string
	fontSize int
	shadow   int
	rotate   int
	fill     bool
}

func newWatermark() Watermark {
	return Watermark{DefaultXMargin, DefaultYMargin, DefaultTransparency, nil}
}

func (img *ImageWand) watermark(o Watermark) (err error) {
	if o.Picture != nil {
		err = o.Picture.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_SET)
		if err != nil {
			return err
		}
		prev := o.Picture.SetImageChannelMask(imagick.CHANNEL_ALPHA)
		err := o.Picture.EvaluateImage(imagick.EVAL_OP_MULTIPLY, o.Transparency)
		if err != nil {
			return err
		}

		err = img.MagickWand.CompositeImage(o.Picture, imagick.COMPOSITE_OP_DISSOLVE, true, o.XMargin, o.YMargin)
		if err != nil {
			return err
		}
		o.Picture.SetImageChannelMask(prev)
	}
	return nil
}

func (img *ImageWand) setTextAsPicture(t Text) (*imagick.MagickWand, error) {
	text := imagick.NewMagickWand()

	img.PixelWand.SetColor(t.color)
	img.DrawWand.SetFillColor(img.PixelWand)

	err := img.DrawWand.SetFont(t.front)
	if err != nil {
		return nil, err
	}
	img.DrawWand.SetFontSize(float64(t.fontSize))
	img.DrawWand.SetGravity(imagick.GRAVITY_CENTER)
	img.DrawWand.Rotate(float64(t.rotate))
	img.DrawWand.Annotation(0, 0, t.text)

	//set text picture
	img.PixelWand.SetColor("none")
	err = text.NewImage(uint(4096), uint(4096), img.PixelWand)
	if err != nil {
		return nil, err
	}
	err = text.DrawImage(img.DrawWand)
	if err != nil {
		return nil, err
	}
	err = text.TrimImage(1)
	if err != nil {
		return nil, err
	}
	return text, nil
}
