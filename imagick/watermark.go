package imagick

import (
	"gopkg.in/gographics/imagick.v3/imagick"
)

const DdfaultGravity = imagick.GRAVITY_SOUTH_EAST
const AlphaOpaque =  imagick.ALPHA_CHANNEL_UNDEFINED
const EvaluateOperator = imagick.EVAL_OP_MULTIPLY

type Watermark struct {
	XMargin      int
	YMargin      int
	Gravity      imagick.GravityType	//used to set putting watermark form where
	Transparency float64
	Picture      *imagick.MagickWand	//watermark is picture
	Text         *Text
}

type Text struct {
	text     string
	textType string
	color    string
	front    string
	fontSize int
	shadow   int
	rotate   int
	fill     bool
}

func newWatermark() Watermark {
	return Watermark{DefaultXMargin, DefaultYMargin, DdfaultGravity, DefaultTransparency, nil, new(Text)}
}

func (img *ImageWand) watermark(o Watermark) (err error) {
	if o.Picture != nil {
		err = o.Picture.SetImageAlphaChannel(AlphaOpaque)
		if err != nil {
			return err
		}
		err = o.Picture.EvaluateImage(EvaluateOperator, o.Transparency)

		if err != nil {
			return err
		}

		err = img.MagickWand.CompositeImage(o.Picture, o.Picture.GetImageCompose(), true, o.XMargin, o.YMargin)
		if err != nil {
			return err
		}
	}
	if o.Text.text != "" {
		img.PixelWand.SetColor(o.Text.color)
		img.DrawWand.SetFillColor(img.PixelWand)

		img.DrawWand.SetFillOpacity(o.Transparency)
		err = img.DrawWand.SetFont(o.Text.front)
		if err != nil {
			return err
		}
		img.DrawWand.SetFontSize(float64(o.Text.fontSize))
		img.DrawWand.SetGravity(o.Gravity)
		img.DrawWand.Annotation(float64(o.XMargin), float64(o.YMargin), o.Text.text)

		err = img.MagickWand.DrawImage(img.DrawWand)
		if err != nil {
			return err
		}
	}
	return nil
}
