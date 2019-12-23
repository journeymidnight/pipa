package imagick

import (
	"github.com/journeymidnight/pipa/helper"
	"github.com/journeymidnight/pipa/library"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func adjustCropTask(plan ResizePlan, width, height uint) {
	//单宽高缩放
	if plan.Width+plan.Height != 0 && plan.Width*plan.Height == 0 {
		return
	}
	//单长短边缩放
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

	//同时指定宽高缩放
	if plan.Width > 0 && plan.Height > 0 {
		if plan.Mode == "lfit" { //长边优先
			if width >= height {
				plan.Height = 0
			} else {
				plan.Width = 0
			}
		}

		if plan.Mode == "mfit" { //短边优先
			if width >= height {
				plan.Width = 0
			} else {
				plan.Height = 0
			}
		}
		return
	}

	//同时指定长短边缩放
	if plan.Long > 0 && plan.Short > 0 {
		if plan.Mode == "lfit" { //长边优先
			if width >= height {
				plan.Width = plan.Long
				plan.Height = 0
			} else {
				plan.Height = plan.Long
				plan.Width = 0
			}
		}

		if plan.Mode == "mfit" { //短边优先
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
	case library.WQYZhenHei:
		return "WQYZH.ttf"
	case library.WQYMicroHei:
		return "WQYWMH.ttf"
	case library.FangZhengShuoSong:
		return "FZSSJW.TTF"
	case library.FangZhengKaiTi:
		return "FZKTJW.TTF"
	case library.FangZhengHeiTi:
		return "FZHTJW.TTF"
	case library.FangZhengFangSong:
		return "FZFSJW.TTF"
	case library.DroidSansFallBack:
		return "DroidSansFallBack.ttf"
	default:
		return "WQYZH.ttf"
	}
}
