package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func CacheMiddleware(maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age="+
			fmt.Sprintf("%d", int(maxAge.Seconds())))
		c.Next()
	}
}
