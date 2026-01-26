package errorx

import "net/http"

// ============ 工具函数 ============

// Is 判断错误是否匹配目标错误码
func Is(err error, target *Error) bool {
	e, ok := err.(*Error)
	return ok && e.Code() == target.Code()
}

// AsError 类型断言，获取 Error 指针
func AsError(err error) (*Error, bool) {
	e, ok := err.(*Error)
	return e, ok
}

// IsRetriable 判断错误是否可重试
func IsRetriable(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Retriable()
	}
	return false
}

// GetCode 从错误中获取错误码，如果不是 Error 返回 0
func GetCode(err error) int {
	if e, ok := err.(*Error); ok {
		return e.Code()
	}
	return 0
}

// GetStatusCode 从错误中获取 HTTP 状态码，如果不是 Error 返回 500
func GetStatusCode(err error) int {
	if e, ok := err.(*Error); ok {
		return e.StatusCode()
	}
	return http.StatusInternalServerError
}

// Cause 递归获取最底层的原始错误
func Cause(err error) error {
	for {
		e, ok := err.(*Error)
		if !ok {
			return err
		}
		if cause := e.Unwrap(); cause != nil {
			err = cause
			continue
		}
		return err
	}
}

// Wrap 用指定的 Error 包装一个错误
func Wrap(baseErr *Error, cause error) *Error {
	return baseErr.With(WithCause(cause))
}

// WrapWithMessage 用指定的 Error 和自定义信息包装一个错误
func WrapWithMessage(baseErr *Error, message string, cause error) *Error {
	return baseErr.With(
		WithMessage(message),
		WithCause(cause),
	)
}
