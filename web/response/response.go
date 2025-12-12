package response

import (
	"net/http"
	"time"
)

// Response 统一响应结构
type Response struct {
	Code      int    `json:"code"`                 // 业务状态码
	Message   string `json:"message"`              // 响应消息
	Data      any    `json:"data,omitempty"`       // 响应数据
	RequestID string `json:"request_id,omitempty"` // 请求追踪ID
	Timestamp int64  `json:"timestamp"`            // 时间戳(毫秒)

	httpStatus int `json:"-"` // HTTP状态码(内部使用)
}

// PageData 分页数据
type PageData struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

// ---------- 预定义业务码 ----------

const (
	CodeSuccess = 0 // 成功
	CodeError   = 1 // 通用错误

	// 客户端错误 4xxxx
	CodeInvalidParams    = 40000 // 参数错误
	CodeUnauthorized     = 40100 // 未授权
	CodeForbidden        = 40300 // 禁止访问
	CodeNotFound         = 40400 // 资源不存在
	CodeMethodNotAllowed = 40500 // 方法不允许
	CodeConflict         = 40900 // 资源冲突
	CodeTooManyRequests  = 42900 // 请求过多

	// 服务端错误 5xxxx
	CodeInternalError      = 50000 // 服务器内部错误
	CodeServiceUnavailable = 50300 // 服务不可用
	CodeGatewayTimeout     = 50400 // 网关超时

	// 业务错误 1xxxx (自定义)
	CodeBusinessError   = 10000 // 业务错误
	CodeRecordExists    = 10001 // 记录已存在
	CodeRecordNotFound  = 10002 // 记录不存在
	CodeOperationFailed = 10003 // 操作失败
)

// ---------- 构造器 ----------

// newResponse 内部构造函数
func newResponse(code int, message string, httpStatus int) *Response {
	return &Response{
		Code:       code,
		Message:    message,
		Timestamp:  time.Now().UnixMilli(),
		httpStatus: httpStatus,
	}
}

// ---------- 成功响应 ----------

// OK 成功响应
func OK() *Response {
	return newResponse(CodeSuccess, "success", http.StatusOK)
}

// Data 成功响应(带数据)
func Data(data any) *Response {
	return OK().WithData(data)
}

// Page 分页响应
func Page(list any, total int64, page, pageSize int) *Response {
	return Data(PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Created 资源创建成功 (201)
func Created(data any) *Response {
	return newResponse(CodeSuccess, "created", http.StatusCreated).WithData(data)
}

// Accepted 请求已接受 (202)
func Accepted(data any) *Response {
	if data == nil {
		return newResponse(CodeSuccess, "accepted", http.StatusAccepted)
	}
	return newResponse(CodeSuccess, "accepted", http.StatusAccepted).WithData(data)
}

// NoContent 无内容响应 (204)
func NoContent() *Response {
	return newResponse(CodeSuccess, "no content", http.StatusNoContent)
}

// ---------- 错误响应 ----------

// Error 通用错误
func Error(message string) *Response {
	return newResponse(CodeError, message, http.StatusOK)
}

// BadRequest 请求参数错误 (400)
func BadRequest(message string) *Response {
	if message == "" {
		message = "invalid parameters"
	}
	return newResponse(CodeInvalidParams, message, http.StatusBadRequest)
}

// Unauthorized 未授权 (401)
func Unauthorized(message string) *Response {
	if message == "" {
		message = "unauthorized"
	}
	return newResponse(CodeUnauthorized, message, http.StatusUnauthorized)
}

// Forbidden 禁止访问 (403)
func Forbidden(message string) *Response {
	if message == "" {
		message = "forbidden"
	}
	return newResponse(CodeForbidden, message, http.StatusForbidden)
}

// NotFound 资源不存在 (404)
func NotFound(message string) *Response {
	if message == "" {
		message = "not found"
	}
	return newResponse(CodeNotFound, message, http.StatusNotFound)
}

// MethodNotAllowed 方法不允许 (405)
func MethodNotAllowed(message string) *Response {
	if message == "" {
		message = "method not allowed"
	}
	return newResponse(CodeMethodNotAllowed, message, http.StatusMethodNotAllowed)
}

// Conflict 资源冲突 (409)
func Conflict(message string) *Response {
	if message == "" {
		message = "conflict"
	}
	return newResponse(CodeConflict, message, http.StatusConflict)
}

// TooManyRequests 请求过于频繁 (429)
func TooManyRequests(message string) *Response {
	if message == "" {
		message = "too many requests"
	}
	return newResponse(CodeTooManyRequests, message, http.StatusTooManyRequests)
}

// InternalError 服务器内部错误 (500)
func InternalError(message string) *Response {
	if message == "" {
		message = "internal server error"
	}
	return newResponse(CodeInternalError, message, http.StatusInternalServerError)
}

// ServiceUnavailable 服务不可用 (503)
func ServiceUnavailable(message string) *Response {
	if message == "" {
		message = "service unavailable"
	}
	return newResponse(CodeServiceUnavailable, message, http.StatusServiceUnavailable)
}

// GatewayTimeout 网关超时 (504)
func GatewayTimeout(message string) *Response {
	if message == "" {
		message = "gateway timeout"
	}
	return newResponse(CodeGatewayTimeout, message, http.StatusGatewayTimeout)
}

// Custom 自定义响应
func Custom(code int, message string, httpStatus int) *Response {
	return newResponse(code, message, httpStatus)
}

// ---------- 链式调用 ----------

// WithData 设置数据
func (r *Response) WithData(data any) *Response {
	r.Data = data
	return r
}

// WithMessage 设置消息
func (r *Response) WithMessage(message string) *Response {
	r.Message = message
	return r
}

// WithRequestID 设置请求ID
func (r *Response) WithRequestID(requestID string) *Response {
	r.RequestID = requestID
	return r
}

// WithCode 设置业务码
func (r *Response) WithCode(code int) *Response {
	r.Code = code
	return r
}

// ---------- 工具方法 ----------

// Status 获取HTTP状态码
func (r *Response) Status() int {
	if r.httpStatus > 0 {
		return r.httpStatus
	}
	return http.StatusOK
}

// IsSuccess 判断是否成功
func (r *Response) IsSuccess() bool {
	return r.Code == CodeSuccess
}

// ---------- Web框架适配 ----------

// JSON 通用适配
// Gin: c.JSON(resp.Status(), resp)
func (r *Response) GJSON(c interface {
	JSON(code int, i any)
}) {
	c.JSON(r.Status(), r)
}

// JSON 通用适配
// Echo: c.JSON(resp.Status(), resp)
func (r *Response) JSON(c interface {
	JSON(code int, i any) error
}) error {
	return c.JSON(r.Status(), r)
}

// Send Fiber 框架适配
// Fiber: return c.Status(resp.Status()).JSON(resp)
func (r *Response) Send(c interface {
	Status(status int) interface{ JSON(data any) error }
}) error {
	return c.Status(r.Status()).JSON(r)
}

// Write 标准库 net/http 适配
func (r *Response) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.Status())

	// 这里需要配合 json.NewEncoder 使用
	// 或者在项目中统一封装 WriteJSON 函数
	return nil
}

// ---------- 快捷方法(业务错误) ----------

// BusinessError 业务错误
func BusinessError(message string) *Response {
	return newResponse(CodeBusinessError, message, http.StatusOK)
}

// RecordExists 记录已存在
func RecordExists(message string) *Response {
	if message == "" {
		message = "record already exists"
	}
	return newResponse(CodeRecordExists, message, http.StatusOK)
}

// RecordNotFound 记录不存在
func RecordNotFound(message string) *Response {
	if message == "" {
		message = "record not found"
	}
	return newResponse(CodeRecordNotFound, message, http.StatusOK)
}

// OperationFailed 操作失败
func OperationFailed(message string) *Response {
	if message == "" {
		message = "operation failed"
	}
	return newResponse(CodeOperationFailed, message, http.StatusOK)
}
