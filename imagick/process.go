package imagick

import (
	. "github.com/journeymidnight/pipa/error"
	"github.com/journeymidnight/pipa/helper"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type ResizePlan struct {
	Mode       string
	Width      int
	Height     int
	Long       int
	Short      int
	Limit      bool
	Color      string
	Proportion int
	Data       []byte
}

type WatermarkPlan struct {
	Transparency int
	Position     string
	XMargin      int
	YMargin      int
	Voffset      int
	PictureMask  WatermarkPicture
	TextMask     WatermarkText
	Order        int //default 0:Image watermark in front 1:Text watermark in front
	Align        int //default 0:align from the top 1:middle 2:bottom
	Interval     int
}

type WatermarkPicture struct {
	Image  string
	Data   []byte
	Rotate RotatePlan
	Crop   CropPlan
}

type WatermarkText struct {
	Text   string
	Type   string
	Color  string
	Size   int
	Shadow int
	Rotate int
	Fill   bool
}

type CropPlan struct {
}

type RotatePlan struct {
}

type ImageWand struct {
	MagickWand *imagick.MagickWand
	PixelWand  *imagick.PixelWand
	DrawWand   *imagick.DrawingWand
}

func Initialize() (lib interface{}) {
	imagick.Initialize()
	return lib
}

func (img *ImageWand) Destory() {
	img.MagickWand.Destroy()
	img.DrawWand.Destroy()
	img.PixelWand.Destroy()
}

func (img *ImageWand) Terminate() {
	imagick.Terminate()
}

func NewImageWand() ImageWand {
	img := ImageWand{
		MagickWand: imagick.NewMagickWand(),
		PixelWand:  imagick.NewPixelWand(),
		DrawWand:   imagick.NewDrawingWand(),
	}
	return img
}

func (img *ImageWand) ResizeImageProcess(data []byte, plan ResizePlan) error {
	helper.Log.Info("start resize image, plan: ", )
	err := img.MagickWand.ReadImageBlob(data)
	if err != nil {
		helper.Log.Error("read data failed", err)
		return err
	}
	originWidth := int(img.MagickWand.GetImageWidth())
	originHeight := int(img.MagickWand.GetImageHeight())
	if originHeight > 30000 || originWidth > 30000 {
		return ErrPictureWidthOrHeightTooLong
	}

	o := newResize()
	o.Limit = plan.Limit
	o.Background = plan.Color

	if plan.Data != nil {
		if plan.Proportion != 0 {
			factor, err := factorCalculations(img, plan.Data, float64(plan.Proportion))
			if err != nil {
				return err
			}
			o.Zoom = factor
			err = img.resize(o)
			if err != nil {
				return err
			}
			return nil
		}
	}

	//proportion zoom
	if plan.Proportion != 0 {
		factor := float64(plan.Proportion) / 100.0
		helper.Log.Info("scaling factor: ", factor)
		o.Zoom = factor
		err = img.resize(o)
		if err != nil {
			return err
		}
		return nil
	}

	o.Width = plan.Width
	o.Height = plan.Height
	switch plan.Mode {
	case "lfit":
		adjustCropTask(plan, originWidth, originHeight)
		helper.Log.Info("trans params ", o)
		err = img.resize(o)
		if err != nil {
			return err
		}
		break
	case "mfit":
		adjustCropTask(plan, originWidth, originHeight)
		helper.Log.Info("trans params ", o)
		err = img.resize(o)
		if err != nil {
			return err
		}
		break
	case "pad":
		o.Pad = true
		helper.Log.Info("trans params ", o)
		err = img.resize(o)
		if err != nil {
			return err
		}
		break
	case "fixed":
		o.Force = true
		helper.Log.Info("trans params ", o)
		err = img.resize(o)
		if err != nil {
			return err
		}
		break
	case "fill":
		o.Crop = true
		helper.Log.Info("trans params ", o)
		err = img.resize(o)
		if err != nil {
			return err
		}
	}
	return nil
}

func (img *ImageWand) ImageWatermarkProcess(data []byte, plan WatermarkPlan) error {
	helper.Log.Info("start resize image, plan: ", )
	err := img.MagickWand.ReadImageBlob(data)
	if err != nil {
		helper.Log.Error("read data failed", err)
		return err
	}
	originWidth := int(img.MagickWand.GetImageWidth())
	originHeight := int(img.MagickWand.GetImageHeight())
	if originHeight > 30000 || originWidth > 30000 {
		return ErrPictureWidthOrHeightTooLong
	}
	w := newWatermark()
	if plan.PictureMask.Image != "" {
		picture := imagick.NewMagickWand()
		err := picture.ReadImageBlob(plan.PictureMask.Data)
		if err != nil {
			helper.Log.Error("open watermark picture file failed")
			return err
		}
		wmWidth := int(picture.GetImageWidth())
		wmHeight := int(picture.GetImageHeight())

		w.Picture = picture
		w.Transparency = plan.Transparency
		switch plan.Position {
		case NorthWest:
			w.XMargin = plan.XMargin
			w.YMargin = plan.YMargin
			break
		case North:
			w.XMargin = (originWidth - wmWidth) / 2
			w.YMargin = plan.YMargin
			break
		case NorthEast:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = plan.YMargin
			break
		case West:
			w.XMargin = plan.XMargin
			w.YMargin = (originHeight-wmHeight)/2 - plan.Voffset
			break
		case Center:
			w.XMargin = (originWidth - wmWidth) / 2
			w.YMargin = (originHeight-wmHeight)/2 - plan.Voffset
			break
		case East:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = (originHeight-wmHeight)/2 - plan.Voffset
			break
		case SouthWest:
			w.XMargin = plan.XMargin
			w.YMargin = originHeight - plan.YMargin - wmHeight
			break
		case South:
			w.XMargin = (originWidth - wmWidth) / 2
			w.YMargin = originHeight - plan.YMargin - wmHeight
			break
		case SouthEast:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = originHeight - plan.YMargin - wmHeight
			break
		default:
			w.XMargin = originWidth - plan.XMargin - wmWidth
			w.YMargin = originHeight - plan.YMargin - wmHeight
		}
		helper.Log.Info("trans params ", w)
		err = img.watermark(w)
		if err != nil {
			return err
		}
		return nil
	} else if plan.TextMask.Text != "" {
		w.Text.text = plan.TextMask.Text
		w.Transparency = plan.Transparency
		w.Text.color = plan.TextMask.Color
		w.Text.textType = selectTextType(plan.TextMask.Type)
		w.Text.fontSize = plan.TextMask.Size
		w.Text.shadow = plan.TextMask.Shadow
		w.Text.rotate = plan.TextMask.Rotate
		w.Text.fill = plan.TextMask.Fill
		switch plan.Position {
		case NorthWest:
			w.Gravity = imagick.GRAVITY_NORTH_WEST
			w.XMargin = plan.XMargin
			w.YMargin = plan.YMargin
			break
		case North:
			w.Gravity = imagick.GRAVITY_NORTH
			w.XMargin = 0
			w.YMargin = plan.YMargin
			break
		case NorthEast:
			w.Gravity = imagick.GRAVITY_NORTH_EAST
			w.XMargin = plan.XMargin
			w.YMargin = plan.YMargin
			break
		case West:
			w.Gravity = imagick.GRAVITY_WEST
			w.XMargin = plan.XMargin
			w.YMargin = -plan.Voffset
			break
		case Center:
			w.Gravity = imagick.GRAVITY_CENTER
			w.XMargin = 0
			w.YMargin = -plan.Voffset
			break
		case East:
			w.Gravity = imagick.GRAVITY_EAST
			w.XMargin = plan.XMargin
			w.YMargin = -plan.Voffset
			break
		case SouthWest:
			w.Gravity = imagick.GRAVITY_SOUTH_WEST
			w.XMargin = plan.XMargin
			w.YMargin = plan.YMargin
			break
		case South:
			w.Gravity = imagick.GRAVITY_SOUTH
			w.XMargin = 0
			w.YMargin = plan.YMargin
			break
		case SouthEast:
			w.Gravity = imagick.GRAVITY_SOUTH_EAST
			w.XMargin = plan.XMargin
			w.YMargin = plan.YMargin
			break
		default:
			w.Gravity = imagick.GRAVITY_SOUTH_EAST
			w.XMargin = plan.XMargin
			w.YMargin = plan.YMargin
		}
		helper.Log.Info("trans params ", w)
		err = img.watermark(w)
		if err != nil {
			helper.Log.Error()
			return err
		}
		return nil
	} else {
		return ErrInvalidWatermarkProcess
	}
}

func (img *ImageWand) ReturnData() []byte {
	return img.MagickWand.GetImageBlob()
}
