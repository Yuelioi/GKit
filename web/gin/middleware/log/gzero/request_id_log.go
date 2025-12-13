package gzero

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type requestIDKey = string

const key requestIDKey = "request_id"

// requestIDLogger 最好先使用 requestid 中间件生成一下
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 gin context 获取 request_id（由 requestid middleware 生成）
		rid := c.GetString("request_id")
		if rid == "" {
			// 如果没有，使用 uuid 生成
			rid = uuid.New().String()
			c.Set("request_id", rid)
		}

		// 将 request_id 注入到 zerolog context
		logger := zerolog.Ctx(c.Request.Context()).With().Str(key, rid).Logger()
		ctx := logger.WithContext(c.Request.Context())

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(key).(string); ok {
		return id
	}
	return ""
}
