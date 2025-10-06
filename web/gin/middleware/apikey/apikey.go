package apikey

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 定义验证函数签名
type Validator func(key string) bool

// ErrorHandler 自定义错误处理函数签名
type ErrorHandler func(c *gin.Context, status int, msg string)

// 配置
type Config struct {
	Header       string       // 请求头名，默认 "Authorization"
	Scheme       string       // 认证前缀，默认 "apikey"
	Validator    Validator    // 必填：验证函数
	ErrorHandler ErrorHandler // 可选：自定义错误响应
}

// 默认错误输出
func defaultErrorHandler(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
	c.Abort()
}

func Default(validator Validator) gin.HandlerFunc {
	return New(Config{
		Validator: validator,
	})
}

// New 创建 API Key 鉴权中间件
func New(cfg Config) gin.HandlerFunc {
	// 默认值处理
	if cfg.Header == "" {
		cfg.Header = "Authorization"
	}
	if cfg.Scheme == "" {
		cfg.Scheme = "apikey"
	}
	if cfg.Validator == nil {
		panic("apikey middleware: Validator is required")
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = defaultErrorHandler
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader(cfg.Header)
		if authHeader == "" {
			cfg.ErrorHandler(c, http.StatusUnauthorized, "未授权")
			return
		}

		// 格式：Authorization: ApiKey xxx
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], cfg.Scheme) {
			cfg.ErrorHandler(c, http.StatusUnauthorized, "无效的认证格式")
			return
		}

		key := parts[1]
		if !cfg.Validator(key) {
			cfg.ErrorHandler(c, http.StatusUnauthorized, "无效的验证信息")
			return
		}

		c.Next()
	}
}
