package library

import "github.com/journeymidnight/pipa/imagick"

const (
	//Resize default param
	Zoom       = 0.0
	Force      = false
	Crop       = false
	Pad        = false
	Limit      = true
	//Watermark default param
	XMargin      = 10
	YMargin      = 10
	Transparency = 100
	FrontSize    = 40.0
	Background   = "#FFFFFF"
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

type Library interface {

}

func NewLibrary() Library {
	//TODO: support other libraries
	return imagick.Initialize()
}