package errorx

// Error 统一错误接口
type Error interface {
	error
	Code() int
	Message() string
	HttpStatus() int
	IsRetriable() bool
	Cause() error
	WithCause(error) Error
}

type customError struct {
	code       int
	message    string
	httpStatus int
	retriable  bool
	cause      error
}

func (e *customError) Error() string     { return e.message }
func (e *customError) Code() int         { return e.code }
func (e *customError) Message() string   { return e.message }
func (e *customError) HttpStatus() int   { return e.httpStatus }
func (e *customError) IsRetriable() bool { return e.retriable }
func (e *customError) Cause() error      { return e.cause }

func (e *customError) WithCause(cause error) Error {
	return &customError{
		code:       e.code,
		message:    e.message,
		httpStatus: e.httpStatus,
		retriable:  e.retriable,
		cause:      cause,
	}
}

// New 创建新错误
func New(code int, message string, httpStatus int) Error {
	return &customError{
		code:       code,
		message:    message,
		httpStatus: httpStatus,
	}
}

// NewRetriable 创建可重试的错误
func NewRetriable(code int, message string, httpStatus int) Error {
	return &customError{
		code:       code,
		message:    message,
		httpStatus: httpStatus,
		retriable:  true,
	}
}
