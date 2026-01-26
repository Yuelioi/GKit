package errorx

// Error 自定义错误结构
type Error struct {
	code       int
	message    string
	httpStatus int
	retriable  bool
	cause      error
}

// 实现 error 接口
func (e *Error) Error() string { return e.message }

// Code 返回错误码
func (e *Error) Code() int { return e.code }

// Message 返回错误信息
func (e *Error) Message() string { return e.message }

// StatusCode 返回 HTTP 状态码
func (e *Error) StatusCode() int { return e.httpStatus }

// Retriable 判断是否可重试
func (e *Error) Retriable() bool { return e.retriable }

// Unwrap 返回底层错误（标准 Go 错误链）
func (e *Error) Unwrap() error { return e.cause }

// ============ 工厂函数 ============

// New 创建新错误
func New(code int, message string, httpStatus int) *Error {
	return &Error{
		code:       code,
		message:    message,
		httpStatus: httpStatus,
		retriable:  false,
	}
}

// NewRetriable 创建可重试的错误
func NewRetriable(code int, message string, httpStatus int) *Error {
	return &Error{
		code:       code,
		message:    message,
		httpStatus: httpStatus,
		retriable:  true,
	}
}

// ============ Option 模式 ============

type Option func(*Error)

// WithMessage 修改错误信息
func WithMessage(message string) Option {
	return func(e *Error) { e.message = message }
}

// WithCause 设置底层错误
func WithCause(cause error) Option {
	return func(e *Error) { e.cause = cause }
}

// With 应用多个选项
func (e *Error) With(opts ...Option) *Error {
	for _, opt := range opts {
		opt(e)
	}
	return e
}
