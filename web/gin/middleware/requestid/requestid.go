package requestid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IDGenerator 定义生成 Request ID 的函数类型
type IDGenerator func() string

// RequestIDConfig 保存配置项
type RequestIDConfig struct {
	headerName string
	contextKey string
	generator  IDGenerator
}

// Option 模式定义
type Option func(*RequestIDConfig)

// WithHeaderName 自定义请求头名
func WithHeaderName(name string) Option {
	return func(cfg *RequestIDConfig) {
		cfg.headerName = name
	}
}

// WithContextKey 自定义 context 中存储的键名
func WithContextKey(key string) Option {
	return func(cfg *RequestIDConfig) {
		cfg.contextKey = key
	}
}

// WithGenerator 自定义 ID 生成函数
func WithGenerator(gen IDGenerator) Option {
	return func(cfg *RequestIDConfig) {
		cfg.generator = gen
	}
}

// RequestID 返回通用的请求 ID 中间件
func RequestID(opts ...Option) gin.HandlerFunc {
	// 默认配置
	cfg := &RequestIDConfig{
		headerName: "X-Request-ID",
		contextKey: "request_id",
		generator:  func() string { return uuid.New().String() },
	}

	// 应用 Option
	for _, opt := range opts {
		opt(cfg)
	}

	return func(c *gin.Context) {
		requestID := c.GetHeader(cfg.headerName)
		if requestID == "" {
			requestID = cfg.generator()
		}

		c.Set(cfg.contextKey, requestID)
		c.Header(cfg.headerName, requestID)

		c.Next()
	}
}
