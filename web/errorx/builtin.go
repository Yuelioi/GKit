package errorx

import "net/http"

// ============ 预定义错误 ============

var (
	// 2xx 成功
	OK = New(0, "Success", http.StatusOK)

	// 4xx 客户端错误
	BadRequest     = New(400001, "Bad Request", http.StatusBadRequest)
	InvalidParams  = New(400001, "Invalid Parameters", http.StatusBadRequest)
	NotFound       = New(400004, "Not Found", http.StatusNotFound)
	MethodNotAllow = New(400005, "Method Not Allowed", http.StatusMethodNotAllowed)
	Conflict       = New(400009, "Conflict", http.StatusConflict)

	// 4xx 认证/授权错误
	Unauthorized = New(401001, "Unauthorized", http.StatusUnauthorized)
	Forbidden    = New(403001, "Forbidden", http.StatusForbidden)

	// 5xx 服务器错误（可重试）
	Internal         = NewRetriable(500001, "Internal Server Error", http.StatusInternalServerError)
	NotImplemented   = New(500002, "Not Implemented", http.StatusNotImplemented)
	ServiceUnavail   = NewRetriable(503001, "Service Unavailable", http.StatusServiceUnavailable)
	DeadlineExceeded = NewRetriable(504001, "Deadline Exceeded", http.StatusGatewayTimeout)

	// 自定义错误
	Unknown = New(999999, "Unknown Error", http.StatusInternalServerError)
)
