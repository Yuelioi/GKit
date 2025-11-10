package model

import (
	"fmt"
	"time"
)

// Response 统一响应结构 0 为成功。
type Response struct {
	Code      int    `json:"code"`                 // 统一状态码 (0为成功，非0为各种业务/系统错误)
	Message   string `json:"message"`              // 响应消息，对用户的友好描述
	Data      any    `json:"data,omitempty"`       // 响应数据
	RequestID string `json:"request_id,omitempty"` // 请求追踪ID (用于链路追踪和日志关联)
	Timestamp int64  `json:"timestamp"`            // 时间戳（毫秒级）
	Error     string `json:"error,omitempty"`      // 详细错误/底层错误堆栈（可选，只用于内部调试，不应透传给用户）
}

// 统一状态码常量
const (
	CodeSuccess = 0 // 成功

	// 1XXXX: 客户端请求/参数错误 (通常对应 HTTP 400 Bad Request)
	CodeInvalidParam   = 10001 // 无效的请求参数
	CodeResourceExist  = 10002 // 资源已存在（冲突）
	CodeMissingParam   = 10003 // 缺少必填参数
	CodeInvalidFormat  = 10004 // 参数格式错误
	CodeTooManyRequest = 10005 // 请求过于频繁（限流）

	// 4XXXX: 鉴权/权限相关错误 (通常对应 HTTP 401/403/404)
	CodeUnauthorized = 40001 // 未授权/Token失效/未登录
	CodeTokenExpired = 40002 // Token 已过期
	CodeForbidden    = 40003 // 权限不足
	CodeNotFound     = 40004 // 资源不存在

	// 9XXXX: 服务器系统/内部错误 (通常对应 HTTP 500/503)
	CodeInternalError   = 90001 // 服务器内部错误
	CodeServiceBusy     = 90002 // 服务繁忙/熔断
	CodeDatabaseError   = 90003 // 数据库操作失败
	CodeCacheError      = 90004 // 缓存操作失败
	CodeRemoteCallError = 90005 // 远程调用失败
)

// ---------- 内部构造函数 ----------

// newResponse 基础构造函数，使用毫秒级时间戳
func newResponse(code int, message string, data any) *Response {
	// 如果消息为空，且不是成功状态，则填充默认错误信息
	if message == "" && code != CodeSuccess {
		message = defaultErrorMessage(code)
	} else if message == "" && code == CodeSuccess {
		// 确保成功时 Message 不为空
		message = "success"
	}

	return &Response{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond), // 毫秒级时间戳
	}
}

// 默认错误提示
func defaultErrorMessage(code int) string {
	switch code {
	case CodeInvalidParam:
		return "请求参数错误或缺失"
	case CodeMissingParam:
		return "缺少必填参数"
	case CodeInvalidFormat:
		return "参数格式不正确"
	case CodeTooManyRequest:
		return "请求过于频繁，请稍后重试"

	case CodeResourceExist:
		return "资源已存在，请勿重复创建"

	case CodeUnauthorized:
		return "请先登录或登录信息已过期"
	case CodeTokenExpired:
		return "登录已过期，请重新登录"
	case CodeForbidden:
		return "权限不足，无法执行该操作"
	case CodeNotFound:
		return "请求的资源不存在"

	case CodeInternalError:
		return "服务器开小差了，请稍后重试"
	case CodeServiceBusy:
		return "系统繁忙，请稍后再试"
	case CodeDatabaseError:
		return "数据库操作失败"
	case CodeCacheError:
		return "缓存操作失败"
	case CodeRemoteCallError:
		return "远程服务调用失败"

	default:
		return "未知错误"
	}
}

// ---------- 成功响应构造函数 (Code=0) ----------

// Success 成功响应（有数据）
func Success(data any) *Response {
	return newResponse(CodeSuccess, "", data)
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(message string, data any) *Response {
	return newResponse(CodeSuccess, message, data)
}

// NoContent 成功响应，但不返回数据 (类似 HTTP 204，但返回统一结构)
func NoContent() *Response {
	return newResponse(CodeSuccess, "no content", nil)
}

// ---------- 错误响应构造函数 (Code!=0) ----------

// Error 通用错误响应
// code: 统一错误码 (非0)
// message: 对用户友好的错误描述
func Error(code int, message string) *Response {
	return newResponse(code, message, nil)
}

// InvalidParam 参数错误响应 (Code=10001)
// message 应该具体说明哪个参数出错
func InvalidParam(message string) *Response {
	return newResponse(CodeInvalidParam, message, nil)
}

// NotFound 资源不存在 (Code=40004)
// resourceName: 可选的资源名称（例如 "用户", "订单"）
func NotFound(resourceName string) *Response {
	msg := defaultErrorMessage(CodeNotFound)
	if resourceName != "" {
		msg = fmt.Sprintf("%s 不存在", resourceName)
	}
	return newResponse(CodeNotFound, msg, nil)
}

// Unauthorized 未登录/鉴权失败 (Code=40001)
func Unauthorized(message string) *Response {
	return newResponse(CodeUnauthorized, message, nil)
}

// InternalError 服务器内部错误 (Code=90001)
// 传入底层 error 对象以便记录详细信息，对外返回通用错误提示
func InternalError(err error) *Response {
	r := newResponse(CodeInternalError, "", nil)
	if err != nil {
		r.Error = err.Error() // 仅在内部/日志中可见
	}
	return r
}

// ---------- 链式调用方法 ----------

// WithRequestID 添加请求ID
func (r *Response) WithRequestID(requestID string) *Response {
	r.RequestID = requestID
	return r
}

// WithTimestamp 更新为当前时间戳
func (r *Response) WithTimestamp() *Response {
	// 更新为毫秒级时间戳
	r.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	return r
}

// WithError 添加详细错误信息（通常用于 InternalError 之后）
func (r *Response) WithError(err error) *Response {
	if err != nil {
		r.Error = err.Error()
	}
	return r
}

// ---------- 便利方法 ----------
func (r *Response) HTTPStatus() int {
	switch {
	case r.Code == CodeSuccess:
		return 200
	case r.Code >= 10000 && r.Code < 20000:
		return 400 // 客户端请求错误
	case r.Code >= 40000 && r.Code < 50000:
		if r.Code == CodeUnauthorized || r.Code == CodeTokenExpired {
			return 401
		}
		if r.Code == CodeForbidden {
			return 403
		}
		if r.Code == CodeNotFound {
			return 404
		}
		return 400
	case r.Code >= 90000:
		return 500
	default:
		return 200
	}
}
