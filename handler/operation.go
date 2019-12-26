package handler

import (
	. "github.com/journeymidnight/pipa/error"
	"github.com/journeymidnight/pipa/imagick"
	"strconv"
	"strings"
)

const (
	RESIZE    = "resize"
	WATERMARK = "watermark"
)

type Operation interface {
	GetType() string
	SetDomain(domain string)
	//whether process watermark picture
	SetIsWatermark(flag bool)
	GetOption(captures map[string]string) (err error)
	GetPictureData(data []byte)
	DoProcess(data []byte) (result []byte, err error)
	Close()
}

type Resize struct {
	isWatermark bool
	img  imagick.ImageWand
	plan imagick.ResizePlan
}

func (r *Resize) GetType() string {
	return RESIZE
}

func (r *Resize) SetDomain(domain string) {

}

func (r *Resize) SetIsWatermark(flag bool) {
	r.isWatermark = flag
}

func (r *Resize) GetOption(captures map[string]string) (err error) {
	if captures["P"] == "" {
		r.plan.WatermarkProportion = 0
	} else {
		r.plan.WatermarkProportion, err = strconv.Atoi(captures["P"])
		if err != nil {
			return err
		}
		if r.plan.WatermarkProportion < 1 || r.plan.WatermarkProportion > 100 {
			return ErrInvalidParameterProportion
		}
	}

	if captures["p"] != "" {
		r.plan.Proportion, err = strconv.Atoi(captures["p"])
		if err != nil {
			return err
		}
		if r.plan.Proportion < 1 || r.plan.Proportion > 1000 {
			return ErrInvalidParameterProportion
		}
	}

	if captures["w"] == "" {
		r.plan.Width = 0
	} else {
		r.plan.Width, err = strconv.Atoi(captures["w"])
		if err != nil {
			return err
		}
		if r.plan.Width < 1 || r.plan.Width > 4096 {
			return ErrInvalidParameterBorder
		}
	}

	if captures["h"] == "" {
		r.plan.Height = 0
	} else {
		r.plan.Height, err = strconv.Atoi(captures["h"])
		if err != nil {
			return err
		}
		if r.plan.Height < 1 || r.plan.Height > 4096 {
			return ErrInvalidParameterBorder
		}
	}

	if captures["l"] == "" {
		r.plan.Long = 0
	} else {
		r.plan.Long, err = strconv.Atoi(captures["l"])
		if err != nil {
			return err
		}
		if r.plan.Long < 1 || r.plan.Long > 4096 {
			return ErrInvalidParameterBorder
		}
	}

	if captures["s"] == "" {
		r.plan.Short = 0
	} else {
		r.plan.Short, err = strconv.Atoi(captures["s"])
		if err != nil {
			return err
		}
		if r.plan.Short < 1 || r.plan.Short > 4096 {
			return ErrInvalidParameterBorder
		}
	}

	if captures["limit"] == "" {
		r.plan.Limit = true
	} else {
		limit, err := strconv.Atoi(captures["limit"])
		if err != nil {
			return err
		} else if limit == 1 {
			r.plan.Limit = true
		} else if limit == 0 {
			r.plan.Limit = false
		} else {
			return ErrInvalidParameterLimit
		}
	}

	r.plan.Color = checkColor(captures["color"])

	if ((r.plan.Width != 0 || r.plan.Height != 0) && (r.plan.Long != 0 || r.plan.Short != 0)) == true {
		return ErrInvalidParameterBorder
	}

	switch captures["m"] {
	case "":
		r.plan.Mode = "lfit"
	case "lfit", "mfit", "fixed":
		r.plan.Mode = captures["m"]
	case "fill", "pad":
		r.plan.Mode = captures["m"]
		if r.plan.Width != 0 && r.plan.Height == 0 {
			r.plan.Height = r.plan.Width
		}
		if r.plan.Width == 0 && r.plan.Height != 0 {
			r.plan.Width = r.plan.Height
		}
	default:
		return ErrInvalidParameterMode
	}
	return nil
}

func (r *Resize) GetPictureData(data []byte) {
	r.plan.Data = data
}

func (r *Resize) DoProcess(data []byte) (result []byte, err error) {
	r.img = imagick.NewImageWand()
	defer r.img.Destory()
	err = r.img.ResizeImageProcess(data, r.plan)
	if err != nil {
		return data, err
	}
	return r.img.ReturnData(), nil
}

func (r *Resize) Close() {
	r.img.Terminate()
}

type Watermark struct {
	domain      string
	isWatermark bool
	img         imagick.ImageWand
	plan        imagick.WatermarkPlan
}

func (w *Watermark) GetType() string {
	return WATERMARK
}

func (w *Watermark) SetDomain(domain string) {
	w.domain = domain
}

func (w *Watermark) SetIsWatermark(flag bool) {
	w.isWatermark = flag
}

func (w *Watermark) GetOption(captures map[string]string) (err error) {
	if w.isWatermark == true {
		return ErrInvalidWatermarkPicture
	}
	if captures["t"] == "" {
		w.plan.Transparency = imagick.Transparency
	} else {
		w.plan.Transparency, _ = strconv.Atoi(captures["t"])
		if w.plan.Transparency < 0 || w.plan.Transparency > 100 {
			return ErrInvalidParameterTransparency
		}
	}

	w.plan.PictureMask.Bucket = captures["bucket"]

	switch captures["g"] {
	case "":
		w.plan.Position = imagick.SouthEast
	case imagick.NorthWest, imagick.North, imagick.NorthEast, imagick.West, imagick.Center,
		imagick.East, imagick.SouthWest, imagick.South, imagick.SouthEast:
		w.plan.Position = captures["g"]
	default:
		return ErrInvalidParameterPosition
	}

	if captures["x"] == "" {
		w.plan.XMargin = 10
	} else {
		w.plan.XMargin, err = strconv.Atoi(captures["x"])
		if err != nil {
			return err
		}
		if w.plan.XMargin < 0 || w.plan.XMargin > 4096 {
			return ErrInvalidParameterXMargin
		}
	}

	if captures["y"] == "" {
		w.plan.YMargin = 10
	} else {
		w.plan.YMargin, err = strconv.Atoi(captures["y"])
		if err != nil {
			return err
		}
		if w.plan.YMargin < 0 || w.plan.YMargin > 4096 {
			return ErrInvalidParameterYMargin
		}
	}

	if captures["voffset"] == "" {
		w.plan.Voffset = 0
	} else {
		w.plan.Voffset, err = strconv.Atoi(captures["voffset"])
		if err != nil {
			return err
		}
		if w.plan.Voffset < -1000 || w.plan.Voffset > 1000 {
			return ErrInvalidParameterVoffset
		}
	}

	if captures["image"] == "" {
		w.plan.PictureMask.Image = ""
	} else {
		w.plan.PictureMask.Image, err = ParseBase64String(captures["image"])
		if err != nil {
			return err
		}
	}

	if captures["text"] == "" {
		w.plan.TextMask.Text = ""
	} else {
		if len(captures["text"]) > 64 {
			return ErrInvalidParameterText
		}
		w.plan.TextMask.Text, err = ParseBase64String(captures["text"])
		if err != nil {
			return err
		}
	}

	if captures["type"] == "" {
		w.plan.TextMask.Type = imagick.DefaultTextType
	} else {
		w.plan.TextMask.Type, err = ParseBase64String(captures["type"])
		if err != nil {
			return err
		}
	}

	w.plan.TextMask.Color = checkColor(captures["color"])

	if captures["size"] == "" {
		w.plan.TextMask.Size = imagick.FrontSize
	} else {
		w.plan.TextMask.Size, err = strconv.Atoi(captures["size"])
		if err != nil {
			return err
		}
		if w.plan.YMargin < 0 || w.plan.YMargin > 1000 {
			return ErrInvalidParameterTextSize
		}
	}

	if captures["rotate"] == "" {
		w.plan.TextMask.Rotate = 0
	} else {
		w.plan.TextMask.Rotate, err = strconv.Atoi(captures["rotate"])
		if err != nil {
			return err
		}
		if w.plan.TextMask.Rotate < 0 || w.plan.TextMask.Rotate > 360 {
			return ErrInvalidParameterRotate
		}
	}

	if captures["fill"] == "" {
		w.plan.TextMask.Fill = false
	} else {
		fill, err := strconv.Atoi(captures["limit"])
		if err != nil {
			return err
		}
		if fill == 1 {
			w.plan.TextMask.Fill = true
		} else if fill == 0 {
			w.plan.TextMask.Fill = false
		} else {
			return ErrInvalidParameterFill
		}
	}

	if captures["order"] == "" {
		w.plan.Order = 0
	} else {
		order, err := strconv.Atoi(captures["order"])
		if err != nil {
			return err
		}
		if order == 1 {
			w.plan.Order = 1
		} else if order == 0 {
			w.plan.Order = 0
		} else {
			return ErrInvalidParameter
		}
	}

	if captures["align"] == "" {
		w.plan.Align = 0
	} else {
		align, err := strconv.Atoi(captures["align"])
		if err != nil {
			return err
		}
		if align == 0 {
			w.plan.Align = 0
		} else if align == 1 {
			w.plan.Align = 1
		} else if align == 2 {
			w.plan.Align = 2
		} else {
			return ErrInvalidParameter
		}
	}

	if captures["interval"] == "" {
		w.plan.Interval = 0
	} else {
		w.plan.Interval, err = strconv.Atoi(captures["interval"])
		if err != nil {
			return err
		}
		if w.plan.Interval < 0 || w.plan.Interval > 1000 {
			return ErrInvalidParameter
		}
	}
	return nil
}

func (w *Watermark) GetPictureData(data []byte) {

}

func (w *Watermark) DoProcess(data []byte) (result []byte, err error) {
	w.img = imagick.NewImageWand()
	defer w.img.Destory()
	if w.plan.PictureMask.Image != "" {
		if w.plan.PictureMask.Bucket != "" {
			domain := strings.Split(w.domain, ".")
			w.domain = UrlHead + w.plan.PictureMask.Bucket + w.domain[len(domain[0]):]
		}
		downloadUrl, operations, err := ParseUrl(w.domain+"/"+w.plan.PictureMask.Image, w.isWatermark)
		if err != nil {
			return nil, err
		}

		w.plan.PictureMask.Data, err = downloadImage(downloadUrl)
		if err != nil {
			return nil, err
		}
		for _, op := range operations {
			op.GetPictureData(data)
			w.plan.PictureMask.Data, err = op.DoProcess(w.plan.PictureMask.Data)
			if err != nil {
				return nil, err
			}
		}
	}

	err = w.img.ImageWatermarkProcess(data, w.plan)
	if err != nil {
		return data, err
	}
	return w.img.ReturnData(), nil
}

func (w *Watermark) Close() {
	w.img.Terminate()
}
