package imagick

import (
	"github.com/journeymidnight/pipa/helper"
	"gopkg.in/gographics/imagick.v3/imagick"
	"math"
)

const Method = imagick.FILTER_LANCZOS

type Resize struct {
	Width            int
	Height           int
	Zoom             float64
	Force            bool
	Crop             bool
	Pad              bool
	LimitEnlargement bool
	Background       string
}

func (img *ImageWand) resize(o *Resize) (err error) {
	originWidth := int(img.MagickWand.GetImageWidth())
	originHeight := int(img.MagickWand.GetImageHeight())

	// image calculations
	factor := imageCalculations(o, originWidth, originHeight)
	helper.Log.Info("resize factor: ", factor)

	if o.LimitEnlargement && !o.Force {
		if originWidth < o.Width || originHeight < o.Height {
			factor = 1.0
			o.Width = originWidth
			o.Height = originHeight
		}
	}
	helper.Log.Info("originWidth", originWidth, " originHeight", originHeight, " factor", factor)
	if o.Zoom != ZoomIsZero {
		err = img.MagickWand.ResizeImage(uint(math.Ceil(float64(originWidth)*o.Zoom)), uint(math.Ceil(float64(originHeight)*o.Zoom)), Method)
		if err != nil {
			helper.Log.Error("MagickWand resize image by zoom failed... err:", err)
			return err
		}
	} else if o.Crop != IsNotCrop {
		err = img.MagickWand.ResizeImage(uint(math.Ceil(float64(originWidth)*factor)), uint(math.Ceil(float64(originHeight)*factor)), Method)
		if err != nil {
			helper.Log.Error("MagickWand resize image by fill failed... err:", err)
			return err
		}
		err = img.cropImage(o, int(math.Ceil(float64(originWidth)*factor)), int(math.Ceil(float64(originHeight)*factor)))
		if err != nil {
			return err
		}
	} else if o.Pad != IsNotPad {
		err = img.MagickWand.ResizeImage(uint(math.Ceil(float64(originWidth)*factor)), uint(math.Ceil(float64(originHeight)*factor)), Method)
		if err != nil {
			helper.Log.Error("MagickWand resize image by pad failed... err:", err)
			return err
		}
		err = img.extentImage(o, int(math.Ceil(float64(originWidth)*factor)), int(math.Ceil(float64(originHeight)*factor)))
		if err != nil {
			return err
		}
	} else if o.Force != IsNotForce {
		err = img.MagickWand.ResizeImage(uint(o.Width), uint(o.Height), Method)
		if err != nil {
			helper.Log.Error("MagickWand resize image by force failed... err:", err)
			return err
		}
	} else {
		err = img.MagickWand.ResizeImage(uint(math.Ceil(float64(originWidth)*factor)), uint(math.Ceil(float64(originHeight)*factor)), Method)
		if err != nil {
			helper.Log.Error("MagickWand resize image failed... err:", err)
			return err
		}
	}

	return nil
}

func newResize() *Resize {
	return &Resize{0, 0, ZoomIsZero, IsNotForce, IsNotCrop, IsNotPad, IsLimitEnlargement, DefaultBackground}
}

func imageCalculations(o *Resize, inWidth, inHeight int) float64 {
	factor := 1.0
	wFactor := float64(o.Width) / float64(inWidth)
	hFactor := float64(o.Height) / float64(inHeight)

	switch {
	case o.Width > 0 && o.Height > 0:
		if o.Crop {
			factor = math.Max(hFactor, wFactor)
		} else {
			factor = math.Min(hFactor, wFactor)
		}
	case o.Width > 0:
		factor = wFactor
	case o.Height > 0:
		factor = hFactor
	// Identity transform
	default:
		o.Width = inWidth
		o.Height = inHeight
		break
	}

	return factor
}

func (img *ImageWand) cropImage(o *Resize, originWidth, originHeight int) error {
	offsetWidth := math.Abs(float64(originWidth-o.Width) / 2)
	offsetHeight := math.Abs(float64(originHeight-o.Height) / 2)
	//CropImage(width,height,x,y)
	err := img.MagickWand.CropImage(uint(o.Width), uint(o.Height), int(offsetWidth), int(offsetHeight))
	if err != nil {
		helper.Log.Error("MagickWand resize crop image... err:", err)
		return err
	}
	return nil
}

func (img *ImageWand) extentImage(o *Resize, originWidth, originHeight int) error {
	img.PixelWand.SetColor(o.Background)
	err := img.MagickWand.SetImageBackgroundColor(img.PixelWand)
	if err != nil {
		helper.Log.Error("MagickWand set image background failed... err:", err)
		return err
	}
	offsetWidth := math.Abs(float64(originWidth-o.Width) / 2)
	offsetHeight := math.Abs(float64(originHeight-o.Height) / 2)
	//ExtentImage(width,height,x,y)
	err = img.MagickWand.ExtentImage(uint(o.Width), uint(o.Height), -int(offsetWidth), -int(offsetHeight))
	if err != nil {
		helper.Log.Error("MagickWand extent image failed... err:", err)
		return err
	}
	return nil
}
