package errorx

import "net/http"

var (
	// 2xx 成功
	OK = New(0, "Success", http.StatusOK)

	// 4xx 客户端错误
	InvalidParams = New(400001, "Invalid Parameters", http.StatusBadRequest)
	NotFound      = New(400004, "Not Found", http.StatusNotFound)
	Unauthorized  = New(401001, "Unauthorized", http.StatusUnauthorized)
	Forbidden     = New(403001, "Forbidden", http.StatusForbidden)

	// 5xx 服务器错误
	Internal         = NewRetriable(500001, "Internal Server Error", http.StatusInternalServerError)
	Unavailable      = NewRetriable(503001, "Service Unavailable", http.StatusServiceUnavailable)
	DeadlineExceeded = NewRetriable(504001, "Deadline Exceeded", http.StatusGatewayTimeout)
)
