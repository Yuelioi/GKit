package errorx

type codeErr struct {
	code       int
	message    string
	httpStatus int
	cause      error
}

func (e *codeErr) Error() string {
	return e.message
}

func (e *codeErr) Code() int {
	return e.code
}

func (e *codeErr) Message() string {
	return e.message
}

func (e *codeErr) HttpStatus() int {
	return e.httpStatus
}

func (e *codeErr) Cause() error {
	return e.cause
}
