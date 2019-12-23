package error

type PipaError int

const (
	Ok PipaError = iota
	ErrInvalidTaskString
	ErrDownloadCode
	StatusRequestEntityTooLarge
	StatusUnsupportedMediaType
	ErrNotFoundOssProcess
	ErrInvalidParameter
	ErrInvalidParameterFormat
	ErrInvalidWatermarkParameter
	ErrInvalidWatermarkPicture
)

type ErrorStruct struct {
	ErrorCode    int
	ErrorMessage string
}

var ErrorCodeResponse = map[PipaError]ErrorStruct{
	Ok: {
		ErrorCode:    200,
		ErrorMessage: "ok",
	},
	ErrInvalidTaskString: {
		ErrorCode:    400,
		ErrorMessage: "Invalid task string from request.",
	},
	ErrDownloadCode: {
		ErrorCode:    401,
		ErrorMessage: "Download response code is not 200",
	},
	StatusRequestEntityTooLarge: {
		ErrorCode:    413,
		ErrorMessage: "Picture too large",
	},
	StatusUnsupportedMediaType: {
		ErrorCode:    415,
		ErrorMessage: "Unsupported Media Type",
	},
	ErrNotFoundOssProcess: {
		ErrorCode:    402,
		ErrorMessage: "Can not parameter x-oss-process.",
	},
	ErrInvalidParameter: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter.",
	},
	ErrInvalidParameterFormat: {
		ErrorCode:    405,
		ErrorMessage: "Invalid parameter.",
	},
	ErrInvalidWatermarkParameter: {
		ErrorCode:    406,
		ErrorMessage: "Invalid parameter.",
	},
	ErrInvalidWatermarkPicture: {
		ErrorCode:    406,
		ErrorMessage: "Invalid watermark picture.",
	},
}

func (e PipaError) ErrorCode() int {
	err, ok := ErrorCodeResponse[e]
	if !ok {
		return 400
	}
	return err.ErrorCode
}

func (e PipaError) Error() string {
	err, ok := ErrorCodeResponse[e]
	if !ok {
		return "We encountered an internal error, please try again."
	}
	return err.ErrorMessage
}
