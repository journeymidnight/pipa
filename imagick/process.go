package imagick

import (
	. "github.com/journeymidnight/pipa/library"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type ImageWand struct {
	MagickWand *imagick.MagickWand
	PixelWand  *imagick.PixelWand
	DrawWand   *imagick.DrawingWand
}

func Initialize() Library {
	imagick.Initialize()
	imageProcess := NewImageWand()
	return imageProcess
}

func (img *ImageWand) Terminate() {
	imagick.Terminate()
}

func NewImageWand() (lib Library) {
	img := ImageWand{
		MagickWand: imagick.NewMagickWand(),
		PixelWand:  imagick.NewPixelWand(),
		DrawWand:   imagick.NewDrawingWand(),
	}
	return &img
}
