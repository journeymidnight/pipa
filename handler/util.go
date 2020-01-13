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
	OssProcess   = "x-oss-process"
	DefaultColor = "#FFFFFF"
)

func ParseUrl(taskUrl string, isWatermark bool) (downloadUrl string, operations []Operation, err error) {
	operations = []Operation{}

	if taskUrl[len(taskUrl)-1:] == "/" {
		taskUrl = taskUrl[:len(taskUrl)-1]
	}

	pos := strings.Index(taskUrl, OssProcess)
	if pos == -1 {
		if isWatermark {
			return taskUrl, operations, nil
		}
		return "", operations, ErrNotFoundOssProcess
	}

	path, err := url.Parse(taskUrl)
	if err != nil {
		return "", operations, err
	}

	isBucketDomain, _ := HasBucketInDomain(path.Hostname(), ".", helper.Config.S3Domain)
	if !isBucketDomain {
		//TODO: if need to support format of url like "s3.test.com/bucketName/objectName",modify it
		return "", operations, ErrIsNotBucketDomain
	}
	//path.Hostname():*.s3.test.com
	host := UrlHead + path.Hostname()
	//path.EscapedPath(): /dir/osstest.jpg
	objectPath := path.EscapedPath()
	downloadUrl = host + objectPath

	processParams := path.Query().Get(OssProcess)
	if "" == processParams {
		return "", operations, ErrNotFoundOssProcess
	}
	params := strings.Split(processParams, "/")
	for _, param := range params[1:] {
		if param == "" {
			continue
		}
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
			break
		case ROTATE:
			if isWatermark {
				return "", operations, ErrInvalidWatermarkRotateParam
			}
			break
		}
		operations = append(operations, operation)
	}

	if len(path.Query()) > 1 {
		downloadUrl = downloadUrl + "?"
		for urlQueryKey, urlQueryValue := range path.Query() {
			if OssProcess != urlQueryKey {
				downloadUrl += urlQueryKey + "=" + urlQueryValue[0] + "&"
			}
		}
		downloadUrl = downloadUrl[:len(downloadUrl)-1]
	}

	return downloadUrl, operations, nil
}

func parseParam(param string) (operation Operation, err error) {
	paramKeys := strings.Split(param, ",")
	switch paramKeys[0] {
	case RESIZE:
		operation = &Resize{}
		captures, err := getKeyAndValue(paramKeys[1:])
		if err != nil {
			return operation, err
		}
		err = operation.SetOption(captures)
		if err != nil {
			return operation, err
		}
		return operation, nil
	case WATERMARK:
		operation = &Watermark{}
		captures, err := getKeyAndValue(paramKeys[1:])
		if err != nil {
			return operation, err
		}
		err = operation.SetOption(captures)
		if err != nil {
			return operation, err
		}
		return operation, nil
	case ROTATE:
		operation = &Rotate{}
		captures := make(map[string]string)
		if len(paramKeys) < 2 {
			return operation, ErrInvalidParameterFormat
		}
		captures[ROTATE] = paramKeys[1]
		err = operation.SetOption(captures)
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
			continue
		}
		keys := strings.Split(param, "_")
		if len(keys) < 2 {
			continue
		}
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
