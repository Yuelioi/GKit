package apikey

import (
	"github.com/gin-gonic/gin"
)

func Example(r *gin.Engine) {

	// 简单用法（推荐）
	r.Use(Default(func(k string) bool {
		return k == "my-secret"
	}))

	// 高级配置（链式）
	r.Use(
		NewBuilder().
			WithHeader("X-API-Key").
			WithScheme("Key").
			WithValidator(func(k string) bool { return k == "12345" }).
			WithErrorHandler(func(c *gin.Context, status int, msg string) {
				c.JSON(status, gin.H{"message": msg})
				c.Abort()
			}).
			Handler(),
	)

}
