package errorx

import "net/http"

// ============ 错误接口 ============

// Code 业务错误接口
type Code interface {
	error
	Code() int
	Message() string
	HttpStatus() int
}

func New(code int, message string, httpStatus int) Code {
	return &codeErr{
		code:       code,
		httpStatus: httpStatus,
		message:    message,
	}
}

// NewWithSpec 根据规范创建错误
// 如果错误码已注册，会自动应用规范中的 message 和 httpStatus
func NewWithSpec(code int, customMessage ...string) Code {
	spec, ok := globalRegistry.codes[code]
	if !ok {
		// 如果没注册，创建时会警告（可选）
		// 但仍然允许创建，灵活性高
		return &codeErr{
			code:       code,
			message:    "unknown error",
			httpStatus: http.StatusInternalServerError,
		}
	}

	msg := spec.Message
	if len(customMessage) > 0 && customMessage[0] != "" {
		msg = customMessage[0]
	}

	return &codeErr{
		code:       code,
		message:    msg,
		httpStatus: spec.HttpStatus,
	}
}

// Wrap 包装原始错误
func Wrap(c Code, cause error) Code {
	return &codeErr{
		code:       c.Code(),
		message:    c.Message(),
		httpStatus: c.HttpStatus(),
		cause:      cause,
	}
}

// WithMessage 修改消息
func WithMessage(c Code, message string) Code {
	ce := c.(*codeErr)
	return &codeErr{
		code:       c.Code(),
		message:    message,
		httpStatus: c.HttpStatus(),
		cause:      ce.cause,
	}
}
