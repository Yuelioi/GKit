package cachecontrol

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Example(r *gin.Engine) {
	// 1️⃣ 默认公共缓存 5 分钟
	r.Use(Default())

	// 2️⃣ 自定义：私有缓存 1 小时
	r.Use(NewBuilder().
		Private().
		WithMaxAge(time.Hour).
		Build())

	// 3️⃣ 禁止缓存
	r.Use(NewBuilder().
		NoStore().
		Build())

	// 4️⃣ 长期静态资源缓存
	r.Use(NewBuilder().
		WithMaxAge(365 * 24 * time.Hour).
		Immutable().
		Build())

}
