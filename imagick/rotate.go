package imagick

import "strconv"

type Rotate struct {
	Degrees    int
	Background string
}

func newRotate() *Rotate {
	return &Rotate{0, DefaultBackground}
}

func (img *ImageWand) rotate(o *Rotate) (err error) {
	img.PixelWand.SetColor(o.Background)
	err = img.MagickWand.RotateImage(img.PixelWand, float64(o.Degrees))
	if err != nil {
		panic(err)
	}
	width := int64(img.MagickWand.GetImageWidth())
	height := int64(img.MagickWand.GetImageHeight())
	repage := strconv.FormatInt(width,10) +  "x" + strconv.FormatInt(height,10) + "+0+0"
	err = img.MagickWand.ResetImagePage(repage)
	if err != nil {
		panic(err)
	}
	return nil
}
