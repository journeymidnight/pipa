package pipa

import (
	. "github.com/journeymidnight/pipa/error"
	"net/url"
	"strconv"
	"strings"
)

const (
	OssProcess   = "?x-oss-process=image/"
	DefaultColor = "#FFFFFF"
)

func parseUrl(taskUrl string) (downloadUrl string, operations []Operation, err error) {
	operations = []Operation{}

	urlFragments := strings.Split(taskUrl, "&")

	pos := strings.Index(urlFragments[0], OssProcess)
	if pos == -1 {
		return "", operations, ErrNotFoundOssProcess
	}

	path, err := url.Parse(urlFragments[0])
	if err != nil {
		return "", operations, err
	}
	host := "http://" + path.Hostname()
	objectPath := path.EscapedPath()
	downloadUrl = host + objectPath
	for i := 1; i < len(urlFragments); i++ {
		downloadUrl += urlFragments[i]
	}

	params := taskUrl[pos+len(OssProcess):]

	for _, param := range strings.Split(params, "/") {
		operation, err := parseParam(param)
		if err != nil {
			return "", operations, err
		}
		if operation.GetType() == WATERMARK {
			operation.SetDomain(host)
		}
		switch operation.GetType() {
		case WATERMARK:
			operation.SetDomain(host)
			operation.SetSecondProcessFlag(true)
		default:
			operation.SetSecondProcessFlag(false)
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
		keys := strings.Split(param, "_")
		if len(keys) > 2 {
			return captures, ErrInvalidParameterFormat
		}
		captures[keys[0]] = keys[1]
	}
	return captures, nil
}

func checkColor(color string) string {
	switch {
	case strings.Contains(color, ","):
		rgb := []int{255, 255, 255}
		colors := strings.Split(color, ",")
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
