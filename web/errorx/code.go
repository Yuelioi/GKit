package errorx

import "net/http"

const (
	CodeOK = 0

	// 4xxxx - 标准 HTTP 错误
	CodeInvalidParams = 40000
	CodeUnauthorized  = 40100
	CodeForbidden     = 40300
	CodeNotFound      = 40400

	// 5xxxx - 标准 HTTP 错误
	CodeInternal    = 50000
	CodeUnavailable = 50300
)

var (
	InvalidParams = &codeErr{
		code:       CodeInvalidParams,
		message:    "invalid parameters",
		httpStatus: http.StatusBadRequest,
	}

	Unauthorized = &codeErr{
		code:       CodeUnauthorized,
		message:    "unauthorized",
		httpStatus: http.StatusUnauthorized,
	}

	Forbidden = &codeErr{
		code:       CodeForbidden,
		message:    "forbidden",
		httpStatus: http.StatusForbidden,
	}

	NotFound = &codeErr{
		code:       CodeNotFound,
		message:    "not found",
		httpStatus: http.StatusNotFound,
	}

	Internal = &codeErr{
		code:       CodeInternal,
		message:    "internal server error",
		httpStatus: http.StatusInternalServerError,
	}

	Unavailable = &codeErr{
		code:       CodeUnavailable,
		message:    "service unavailable",
		httpStatus: http.StatusServiceUnavailable,
	}
)
