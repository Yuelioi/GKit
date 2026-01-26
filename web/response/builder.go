package response

import (
	"time"

	"github.com/Yuelioi/gkit/web/errorx"
)

func Success(data interface{}) *Response {
	return &Response{
		Code:       0,
		Message:    "Success",
		Data:       data,
		Timestamp:  time.Now().UnixMilli(),
		httpStatus: 200,
	}
}

// Error 错误响应
func Error(err error) *Response {
	if err == nil {
		return Success(nil)
	}

	// 如果是自定义错误，使用错误信息
	if e, ok := err.(errorx.Error); ok {
		return &Response{
			Code:       e.Code(),
			Message:    e.Message(),
			Timestamp:  time.Now().UnixMilli(),
			httpStatus: e.HttpStatus(),
		}
	}

	// 其他错误当作内部错误处理
	return &Response{
		Code:       errorx.Internal.Code(),
		Message:    errorx.Internal.Message(),
		Timestamp:  time.Now().UnixMilli(),
		httpStatus: errorx.Internal.HttpStatus(),
	}
}

// ============ Builder 链式调用 ============

func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

func (r *Response) WithMessage(msg string) *Response {
	r.Message = msg
	return r
}

func (r *Response) WithTraceID(traceID string) *Response {
	r.TraceID = traceID
	return r
}

func (r *Response) WithCode(code int) *Response {
	r.Code = code
	return r
}

func (r *Response) WithStatus(status int) *Response {
	r.httpStatus = status
	return r
}

// ============ 获取 HTTP 状态码 ============

func (r *Response) Status() int {
	if r.httpStatus > 0 {
		return r.httpStatus
	}
	return 200
}
