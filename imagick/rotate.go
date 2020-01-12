package imagick

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
	return nil
}
