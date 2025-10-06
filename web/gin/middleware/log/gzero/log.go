package gzero

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type GinLoggerBuilder struct {
	logger      zerolog.Logger
	ignorePaths map[string]struct{}
}

func NewGinLoggerBuilder(logger zerolog.Logger) *GinLoggerBuilder {
	return &GinLoggerBuilder{
		logger:      logger,
		ignorePaths: make(map[string]struct{}),
	}
}

// 添加要忽略的路径
func (b *GinLoggerBuilder) IgnorePath(path string) *GinLoggerBuilder {
	b.ignorePaths[path] = struct{}{}
	return b
}

func (b *GinLoggerBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if _, skip := b.ignorePaths[path]; skip {
			c.Next()
			return
		}

		raw := c.Request.URL.RawQuery
		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if errorMessage == "" && statusCode >= 400 {
			errorMessage = fmt.Sprintf("HTTP %d", statusCode)
		}

		fullPath := path
		if raw != "" {
			fullPath = path + "?" + raw
		}

		var evt *zerolog.Event
		switch {
		case statusCode >= 500:
			evt = b.logger.Error()
		case statusCode >= 400:
			evt = b.logger.Warn()
		case statusCode >= 300:
			evt = b.logger.Info()
		default:
			evt = b.logger.Debug()
		}

		evt.Str("client_ip", clientIP).
			Str("method", method).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("path", fullPath).
			Str("user_agent", c.Request.UserAgent()).
			Str("error", errorMessage).
			Msg("HTTP request")
	}
}
