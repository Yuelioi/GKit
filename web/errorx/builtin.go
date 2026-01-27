package errorx

import "net/http"

var (
	// ======================
	// 2xx 成功
	// ======================
	OK = New(0, "Success", http.StatusOK)

	// ======================
	// 4xx 客户端错误
	// ======================

	// 400 Bad Request
	BadRequest    = New(400001, "Invalid Parameters", http.StatusBadRequest)
	MissingParams = New(400002, "Missing Required Parameters", http.StatusBadRequest)
	InvalidFormat = New(400003, "Invalid Data Format", http.StatusBadRequest)

	// 401 / 403 认证 & 授权
	Unauthorized = New(401001, "Unauthorized", http.StatusUnauthorized)
	Forbidden    = New(403001, "Forbidden", http.StatusForbidden)

	// 404 资源不存在
	NotFound = New(404001, "Resource Not Found", http.StatusNotFound)

	// 405 方法不允许
	MethodNotAllow = New(405001, "Method Not Allowed", http.StatusMethodNotAllowed)

	// ======================
	// 409 冲突
	// ======================
	Conflict      = New(409001, "Conflict", http.StatusConflict)
	DuplicateData = New(409002, "Data Already Exists", http.StatusConflict)
	InvalidState  = New(409003, "Invalid Resource State", http.StatusConflict)

	// ======================
	// 422 业务校验失败
	// ======================
	ValidationFailed = New(422001, "Validation Failed", http.StatusUnprocessableEntity)
	ConstraintError  = New(422002, "Constraint Violated", http.StatusUnprocessableEntity)

	// ======================
	// 429 限流
	// ======================
	TooManyRequests = NewRetriable(429001, "Rate Limited", http.StatusTooManyRequests)

	// ======================
	// 5xx 服务器错误（可重试）
	// ======================
	Internal       = NewRetriable(500001, "Internal Server Error", http.StatusInternalServerError)
	NotImplemented = New(501001, "Not Implemented", http.StatusNotImplemented)
	BadGateway     = NewRetriable(502001, "Bad Gateway", http.StatusBadGateway)
	ServiceUnavail = NewRetriable(503001, "Service Unavailable", http.StatusServiceUnavailable)
	Timeout        = NewRetriable(504001, "Request Timeout", http.StatusGatewayTimeout)

	// 第三方服务错误
	ExternalError = NewRetriable(502002, "External Service Error", http.StatusBadGateway)

	// ======================
	// 未知错误
	// ======================
	Unknown = New(999999, "Unknown Error", http.StatusInternalServerError)
)
