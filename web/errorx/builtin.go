package errorx

import "net/http"

const Version = "v1"

const (
	CodeOK            = 0
	CodeInvalidParams = 40000
	CodeUnauthorized  = 40100
	CodeForbidden     = 40300
	CodeNotFound      = 40400
	CodeInternal      = 50000
	CodeUnavailable   = 50300
)

func init() {
	RegisterMust(CodeSpec{
		Code:       CodeOK,
		MessageKey: "error.success",
		HttpStatus: http.StatusOK,
		Version:    Version,
	})

	RegisterMust(CodeSpec{
		Code:       CodeInvalidParams,
		MessageKey: "error.invalid_params",
		HttpStatus: http.StatusBadRequest,
		Version:    Version,
	})

	RegisterMust(CodeSpec{
		Code:       CodeUnauthorized,
		MessageKey: "error.unauthorized",
		HttpStatus: http.StatusUnauthorized,
		Version:    Version,
	})

	RegisterMust(CodeSpec{
		Code:       CodeInternal,
		MessageKey: "error.internal",
		HttpStatus: http.StatusInternalServerError,
		Version:    Version,
		Retriable:  true,
	})
}
