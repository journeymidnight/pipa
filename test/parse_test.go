package test

import (
	"github.com/journeymidnight/pipa/handle"
	"testing"
)

const (
	BucketDomain    = "http://myBucket.s3.test.com/"
	ObjectName      = "ossImages.jpg"
	Key             = "?x-oss-process=image/"
	Resize          = "resize,m_"
	Width           = "w_"
	Height          = "h_"
	Watermark       = "watermark,image_aW1hZ2VzLmpwZz94LW9zcy1wcm9jZXNzPWltYWdlL3Jlc2l6ZSxQXzI5"
	Position        = "g_" //nw,north,ne,west,ceter,east,sw,south,se
	XMargin         = "x_"
	YMargin         = "y_"
	WrongParamOne   = "dsafsagf_"
	WrongParamTwo   = "_fdgsdgdsagf"
	WrongParamThree = "dffad_dgfsaf"
	OtherParams      = "&OSSAccessKeyId=TMP.hjr6bCq2SP9CWFsTUDrFaamfDBh3fT2G6SEZrF1NJE7JTMRTqRkQrNP7e84oYMuoXedhYKcmG6Za4q7YTKmS9PaHeAzUpeZUz62HxJ3Fq2KnzTNct8PHAV8JftAh3G.tmp" +
		"&Signature=Ndgi%2Bzgp3j6hQdVs4WRiWbmUdJ0%3D"
)

func Test_ParseUrl(t *testing.T) {
	urlRightOne := BucketDomain + ObjectName + Key + Resize + "lfit" + Width + "150" +
		Height + "150" + Watermark + Position + "center" + XMargin + "140" + YMargin + "150"
	downloadUrl, _, err := handle.ParseUrl(urlRightOne)
	if err != nil {
		t.Fatal("ParseUrl is wrong")
	} else if downloadUrl != BucketDomain+ObjectName {
		t.Fatal("downloadUrl is wrong")
	}

	urlRightTwo := BucketDomain + ObjectName + Key + Resize + "lfit" + Width + "150" +
		Height + "150" + Watermark + Position + "center" + XMargin + "140" + "y"
	downloadUrl, _, err = handle.ParseUrl(urlRightTwo)
	if err != nil {
		t.Fatal("ParseUrl is wrong")
	} else if downloadUrl != BucketDomain+ObjectName {
		t.Fatal("downloadUrl is wrong")
	}

	urlRightThree := BucketDomain + ObjectName + Key + Resize + "lfit" + Width + "150" +
		Height + "150" + Watermark + Position + "center" + XMargin + "140" + "y" + OtherParams
	downloadUrl, _, err = handle.ParseUrl(urlRightThree)
	if err != nil {
		t.Fatal("ParseUrl is wrong")
	} else if downloadUrl != BucketDomain+ObjectName+OtherParams {
		t.Fatal("downloadUrl is wrong")
	}


}
