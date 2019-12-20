package error

type PipaError int

const (
	Ok PipaError = iota
	ErrInvalidTaskString
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
