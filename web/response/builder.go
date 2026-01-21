package response

import (
	"time"

	"github.com/Yuelioi/gkit/web/errorx"
	"github.com/Yuelioi/gkit/web/i18n"
)

const APIVersion = "v1"

func OK(locale i18n.Locale) *Response {
	return FromCode(errorx.CodeOK, locale)
}

func FromError(err error, locale i18n.Locale) *Response {
	if err == nil {
		return OK(locale)
	}

	if c, ok := err.(errorx.Code); ok {
		return FromCode(c.Code(), locale)
	}

	return FromCode(errorx.CodeInternal, locale)
}

func FromCode(code int, locale i18n.Locale) *Response {
	spec, ok := errorx.GetSpec(code)
	if !ok {
		spec = errorx.CodeSpec{
			Code:       code,
			MessageKey: "error.internal",
			HttpStatus: 500,
		}
	}

	return &Response{
		Code:       spec.Code,
		Message:    i18n.Translate(i18n.Key(spec.MessageKey), locale),
		Timestamp:  time.Now().UnixMilli(),
		Version:    APIVersion,
		httpStatus: spec.HttpStatus,
	}
}
