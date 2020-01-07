package imagick

import (
	. "github.com/journeymidnight/pipa/error"
	"github.com/journeymidnight/pipa/helper"
	"gopkg.in/gographics/imagick.v3/imagick"
)

const (
	//Resize default param
	ZoomIsZero         = 0.0
	IsNotForce         = false
	IsNotCrop          = false
	IsNotPad           = false
	IsLimitEnlargement = true
	//Watermark default param
	DefaultXMargin      = 10
	DefaultYMargin      = 10
	DefaultTransparency = 100
	DefaultFrontSize    = 40.0
	DefaultBackground   = "#FFFFFF"
	DefaultTextColor    = "#000000"
)

const (
	NorthWest = "nw"
	North     = "north"
	NorthEast = "ne"
	West      = "west"
	Center    = "center"
	East      = "east"
	SouthWest = "sw"
	South     = "south"
	SouthEast = "se"
)

//Text type
const (
	DefaultTextType   = "wqy-zenhei"
	WQYZhenHei        = "wqy-zenhei"
	WQYMicroHei       = "wqy-microhei"
	FangZhengShuoSong = "fangzhengshusong"
	FangZhengKaiTi    = "fangzhengkaiti"
	FangZhengHeiTi    = "fangzhengheiti"
	FangZhengFangSong = "fangzhengfangsong"
	DroidSansFallBack = "droidsansfallback"
)

func adjustCropTask(plan *ResizePlan, width, height int) {
	//resize by width or height
	if plan.Width+plan.Height != 0 && plan.Width*plan.Height == 0 {
		return
	}
	//resize by long or short
	if plan.Long+plan.Short != 0 && plan.Long*plan.Short == 0 {
		if plan.Long != 0 {
			if width >= height {
				plan.Width = plan.Long
				plan.Height = 0
			} else {
				plan.Height = plan.Long
				plan.Width = 0
			}
		} else {
			if width >= height {
				plan.Height = plan.Short
				plan.Width = 0
			} else {
				plan.Width = plan.Short
				plan.Height = 0
			}
		}
		return
	}

	//resize by width and height
	if plan.Width > 0 && plan.Height > 0 {
		if plan.Mode == "lfit" { //long first
			if width >= height {
				plan.Height = 0
			} else {
				plan.Width = 0
			}
		}

		if plan.Mode == "mfit" { //short first
			if width >= height {
				plan.Width = 0
			} else {
				plan.Height = 0
			}
		}
		return
	}

	//resize by long and short
	if plan.Long > 0 && plan.Short > 0 {
		if plan.Mode == "lfit" { //long first
			if width >= height {
				plan.Width = plan.Long
				plan.Height = 0
			} else {
				plan.Height = plan.Long
				plan.Width = 0
			}
		}

		if plan.Mode == "mfit" {
			if width >= height { //short first
				plan.Height = plan.Short
				plan.Width = 0
			} else {
				plan.Width = plan.Short
				plan.Height = 0
			}
		}
		return
	}
	return
}

func factorCalculations(watermarkPicture *ImageWand, originPicture []byte, factor float64) (float64, error) {
	picture := imagick.NewMagickWand()
	defer picture.Destroy()
	err := picture.ReadImageBlob(originPicture)
	if err != nil {
		helper.Log.Error("open origin picture file failed")
		return 0, err
	}
	originWidth := float64(picture.GetImageWidth())
	originHeight := float64(picture.GetImageHeight())
	width := float64(watermarkPicture.MagickWand.GetImageWidth())
	height := float64(watermarkPicture.MagickWand.GetImageHeight())
	factor = float64(factor) / 100.0
	tempWidth := originWidth * factor
	tempHeight := originHeight * factor
	widthFactor := tempWidth / width
	heightFactor := tempHeight / height
	if widthFactor > heightFactor {
		factor = heightFactor
	} else {
		factor = widthFactor
	}
	return factor, nil
}

func selectTextType(tType string) string {
	switch tType {
	case WQYZhenHei:
		return "WQYZH.ttf"
	case WQYMicroHei:
		return "WQYWMH.ttf"
	case FangZhengShuoSong:
		return "FZSSJW.TTF"
	case FangZhengKaiTi:
		return "FZKTJW.TTF"
	case FangZhengHeiTi:
		return "FZHTJW.TTF"
	case FangZhengFangSong:
		return "FZFSJW.TTF"
	case DroidSansFallBack:
		return "DroidSansFallBack.ttf"
	default:
		return "WQYZH.ttf"
	}
}

func originPictureIsIllegal(originWidth, originHeight int) error {
	if originHeight > 30000 || originWidth > 30000 {
		return ErrPictureWidthOrHeightTooLong
	}
	if originHeight <= 0 || originWidth <= 0 {
		return ErrPictureWidthOrHeightIsZero
	}
	return nil
}
