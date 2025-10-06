package apikey

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 定义验证函数签名
type Validator func(key string) bool

// 自定义错误处理函数签名
type ErrorHandler func(c *gin.Context, status int, msg string)

type Builder struct {
	header       string
	scheme       string
	validator    Validator
	errorHandler ErrorHandler
}

// 默认错误输出
func defaultErrorHandler(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
	c.Abort()
}

func NewBuilder() *Builder {
	return &Builder{
		header:       "Authorization",
		scheme:       "apikey",
		errorHandler: defaultErrorHandler,
	}
}

func Default(validator Validator) gin.HandlerFunc {
	return NewBuilder().WithValidator(validator).Handler()
}

// WithHeader 设置认证头名
func (b *Builder) WithHeader(header string) *Builder {
	b.header = header
	return b
}

// WithScheme 设置认证前缀（例如 "apikey"）
func (b *Builder) WithScheme(scheme string) *Builder {
	b.scheme = scheme
	return b
}

// WithValidator 设置验证逻辑（必填）
func (b *Builder) WithValidator(fn Validator) *Builder {
	b.validator = fn
	return b
}

// WithErrorHandler 设置自定义错误处理
func (b *Builder) WithErrorHandler(fn ErrorHandler) *Builder {
	b.errorHandler = fn
	return b
}

// Handler 构建 gin.HandlerFunc
func (b *Builder) Handler() gin.HandlerFunc {
	if b.validator == nil {
		panic("apikey middleware: Validator is required")
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader(b.header)
		if authHeader == "" {
			b.errorHandler(c, http.StatusUnauthorized, "未授权")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], b.scheme) {
			b.errorHandler(c, http.StatusUnauthorized, "无效的认证格式")
			return
		}

		key := parts[1]
		if !b.validator(key) {
			b.errorHandler(c, http.StatusUnauthorized, "无效的验证信息")
			return
		}

		c.Next()
	}
}
