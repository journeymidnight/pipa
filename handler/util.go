package handler

import (
	"encoding/base64"
	. "github.com/journeymidnight/pipa/error"
	"github.com/journeymidnight/pipa/helper"
	"net/url"
	"strconv"
	"strings"
)

const (
	UrlHead      = "http://"
	OssProcess   = "?x-oss-process=image/"
	DefaultColor = "#000000"
)

func ParseUrl(taskUrl string, isWatermark bool) (downloadUrl string, operations []Operation, err error) {
	operations = []Operation{}

	if taskUrl[len(taskUrl)-1:] == "/" {
		taskUrl = taskUrl[:len(taskUrl)-1]
	}

	urlFragments := strings.Split(taskUrl, "\u0026")

	pos := strings.Index(urlFragments[0], OssProcess)
	if pos == -1 {
		if isWatermark {
			return taskUrl, operations, nil
		}
		return "", operations, ErrNotFoundOssProcess
	}

	path, err := url.Parse(urlFragments[0])
	if err != nil {
		return "", operations, err
	}
	isBucketDomain, _ := HasBucketInDomain(path.Hostname(), ".", helper.Config.S3Domain)
	if !isBucketDomain {
		return "", operations, ErrIsNotBucketDomain
	}
	//path.Hostname():*.s3.test.com
	host := UrlHead + path.Hostname()
	// /osstest.jpg
	objectPath := path.EscapedPath()
	downloadUrl = host + objectPath
	if len(urlFragments) > 1 {
		downloadUrl = downloadUrl + "?" + urlFragments[1]
		for i := 2; i < len(urlFragments); i++ {
			downloadUrl += "&" + urlFragments[i]
		}
	}

	params := urlFragments[0][pos+len(OssProcess):]
	for _, param := range strings.Split(params, "/") {
		operation, err := parseParam(param)
		if err != nil {
			return "", operations, err
		}
		switch operation.GetType() {
		case WATERMARK:
			if isWatermark {
				return "", operations, ErrWatermarkCanNotProcess
			}
			operation.SetDomain(host)
			operation.SetIsWatermark(true)
		}
		operations = append(operations, operation)
	}
	return downloadUrl, operations, nil
}

func parseParam(param string) (operation Operation, err error) {
	paramKeys := strings.Split(param, ",")

	captures, err := getKeyAndValue(paramKeys[1:])
	if err != nil {
		return operation, err
	}

	switch paramKeys[0] {
	case RESIZE:
		operation = &Resize{}
		err = operation.GetOption(captures)
		if err != nil {
			return operation, err
		}
		return operation, nil
	case WATERMARK:
		operation = &Watermark{}
		err = operation.GetOption(captures)
		if err != nil {
			return operation, err
		}
		return operation, nil
	default:
		return operation, ErrInvalidParameter
	}
}

func getKeyAndValue(paramKeys []string) (captures map[string]string, err error) {
	captures = make(map[string]string)
	for _, param := range paramKeys {
		if param == "" {
			return captures, ErrInvalidParametersHaveSpaces
		}
		keys := strings.Split(param, "_")
		//if len(keys) > 2 {
		//	return captures, ErrInvalidParameterFormat
		//}
		captures[keys[0]] = param[len(keys[0])+1:]
	}
	return captures, nil
}

func checkColor(color string) string {
	switch {
	case strings.Contains(color, "-"):
		rgb := []int{255, 255, 255}
		colors := strings.Split(color, "-")
		for i, num := range colors {
			n, err := strconv.Atoi(num)
			if err != nil {
				break
			}
			if n > 255 {
				continue
			}
			rgb[i] = n
		}
		return "rgb(" + string(rgb[0]) + "," + string(rgb[1]) + "," + string(rgb[2]) + ")"
	case len(color) == 6:
		return "#" + color
	default:
		return DefaultColor
	}
}

func ParseBase64String(str string) (string, error) {
	//str = strings.Replace(str, "_", "/", -1)
	//str = strings.Replace(str, "-", "+", -1)
	mod4String := len(str) % 4
	equalSign := []string{"", "===", "==", "="}
	str += equalSign[mod4String]

	//data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(str))
	data, err := base64.URLEncoding.DecodeString(strings.TrimSpace(str))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func HasBucketInDomain(host string, prefix string, domains []string) (ok bool, bucket string) {
	for _, d := range domains {
		if strings.HasSuffix(host, prefix+d) {
			return true, strings.TrimSuffix(host, prefix+d)
		}
	}
	return false, ""
}
