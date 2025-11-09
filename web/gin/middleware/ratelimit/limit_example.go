package ratelimit

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 限流策略说明：
//
// 1. 全局限流 (WithGlobalLimit)
//    - 作用：限制服务器总体请求量，防止过载
//    - 适用：所有请求
//    - 示例：WithGlobalLimit(1000, time.Second) // 每秒最多 1000 次请求
//
// 2. IP 限流 (WithIPLimit)
//    - 作用：限制单个 IP 的请求量，防止滥用
//    - 适用：所有请求，按 IP 独立计数
//    - 示例：WithIPLimit(100, time.Second) // 每个 IP 每秒最多 100 次请求
//
// 3. 方法限流 (WithMethodLimit) 每个IP独立计算
//    - 作用：限制特定 HTTP 方法的请求量
//    - 适用：POST、PUT、DELETE 等写操作
//    - 示例：WithMethodLimit([]string{"POST", "DELETE"}, 10, time.Second)
//
// 4. 路由限流 (WithRouteLimit) 每个IP独立计算
//    - 作用：限制特定路由的请求量，支持指定 HTTP 方法
//    - 适用：敏感操作，如上传、登录、注册
//    - 示例：
//      WithRouteLimit("/api/v1/files/upload", 1, time.Second, "POST")  // 只限制 POST
//      WithRouteLimit("/api/v1/admin/settings", 10, time.Second, "*")  // 限制所有方法
//      WithRouteLimit("/api/v1/admin/*", 10, time.Second, "POST", "DELETE") // 只限制 POST 和 DELETE
//    - 如果不指定方法参数或传 "*"，则限制所有 HTTP 方法
//
// 5. 优先级顺序（按检查顺序）
//    全局限流 -> 方法限流 -> 路由限流 -> IP 限流
//    任何一个限流触发都会拒绝请求
//
// 6. 推荐配置
//    - 公开 API：全局 + IP 限流
//    - 写操作：+ 方法限流
//    - 敏感操作：+ 路由限流
//    - 清理间隔：1-5 分钟

// Example1_Basic 基础使用示例
func Example1_Basic() {
	r := gin.New()

	// 默认限流：全局 100 次/秒，每个 IP 10 次/秒
	r.Use(Default())

	r.GET("/api/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	r.Run(":8080")
}

// Example2_CustomIPLimit 自定义 IP 限流
func Example2_CustomIPLimit() {
	r := gin.New()

	// 每个 IP 每秒最多 120 次请求
	r.Use(NewBuilder().
		WithIPLimit(120, time.Second).
		WithCleanup(1 * time.Minute).
		Middleware())

	r.GET("/api/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	r.Run(":8080")
}

// Example3_GlobalAndIPLimit 全局 + IP 限流
func Example3_GlobalAndIPLimit() {
	r := gin.New()

	// 服务器总请求限流 + 每个 IP 限流
	r.Use(NewBuilder().
		WithGlobalLimit(1000, time.Second). // 服务器每秒最多 1000 次请求
		WithIPLimit(50, time.Second).       // 每个 IP 每秒最多 50 次请求
		WithCleanup(1 * time.Minute).
		Middleware())

	r.GET("/api/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	r.Run(":8080")
}

// Example4_MethodLimit 方法级别限流（推荐）
func Example4_MethodLimit() {
	r := gin.New()

	// 只限制写操作（POST/PUT/DELETE），读操作不限制
	r.Use(NewBuilder().
		WithMethodLimit([]string{"POST", "PUT", "DELETE"}, 5, time.Second). // 写操作每秒 5 次
		WithCleanup(1 * time.Minute).
		Middleware())

	r.POST("/api/create", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "created"})
	})

	r.GET("/api/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"}) // 不受限制
	})

	r.Run(":8080")
}

// Example5_RouteLimit 特定路由限流
func Example5_RouteLimit() {
	r := gin.New()

	// 对上传接口进行严格限流，指定 HTTP 方法
	r.Use(NewBuilder().
		WithRouteLimit("/api/v1/files/upload", 1, time.Second, "POST").  // 只限制 POST 请求
		WithRouteLimit("/api/v1/images/upload", 1, time.Second, "POST"). // 只限制 POST 请求
		WithRouteLimit("/api/v1/users/profile", 5, time.Second, "PUT").  // 只限制 PUT 请求
		WithRouteLimit("/api/v1/admin/settings", 10, time.Second, "*").  // 限制所有方法
		WithCleanup(1 * time.Minute).
		Middleware())

	r.POST("/api/v1/files/upload", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "file uploaded"}) // 受限：每秒 1 次
	})

	r.POST("/api/v1/images/upload", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "image uploaded"}) // 受限：每秒 1 次
	})

	r.GET("/api/v1/files/upload", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "upload form"}) // 不受限（只限制了 POST）
	})

	r.PUT("/api/v1/users/profile", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "profile updated"}) // 受限：每秒 5 次
	})

	r.GET("/api/v1/admin/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "settings"}) // 受限：每秒 10 次（所有方法）
	})

	r.Run(":8080")
}

// Example6_RoutePrefixLimit 路由前缀匹配限流
func Example6_RoutePrefixLimit() {
	r := gin.New()

	// 限制所有 /api/v1/admin/ 下的请求
	r.Use(NewBuilder().
		WithRouteLimit("/api/v1/admin/*", 10, time.Second, "*").          // 所有方法
		WithRouteLimit("/api/v1/write/*", 5, time.Second, "POST", "PUT"). // 只限制写操作
		WithCleanup(1 * time.Minute).
		Middleware())

	r.GET("/api/v1/admin/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "users"}) // 受限：每秒 10 次
	})

	r.POST("/api/v1/admin/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "settings"}) // 受限：每秒 10 次
	})

	r.POST("/api/v1/write/article", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "article created"}) // 受限：每秒 5 次
	})

	r.GET("/api/v1/write/article", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "article list"}) // 不受限（只限制了 POST/PUT）
	})

	r.GET("/api/v1/public/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "data"}) // 不受限
	})

	r.Run(":8080")
}

// Example7_CombinedLimits 组合限流（推荐生产环境配置）
func Example7_CombinedLimits() {
	r := gin.New()

	// 创建一个共享的限流器实例，配置多种限流策略
	limiter := NewBuilder().
		WithGlobalLimit(1000, time.Second).                                  // 服务器总请求：每秒 1000 次
		WithIPLimit(100, time.Second).                                       // 每个 IP：每秒 100 次
		WithMethodLimit([]string{"POST", "PUT", "DELETE"}, 10, time.Second). // 写操作：每秒 10 次
		WithRouteLimit("/api/v1/files/upload", 1, time.Second, "POST").      // 文件上传：每秒 1 次（仅 POST）
		WithRouteLimit("/api/v1/images/upload", 1, time.Second, "POST").     // 图片上传：每秒 1 次（仅 POST）
		WithCleanup(1 * time.Minute)

	// 全局使用
	r.Use(limiter.Middleware())

	// 定义路由
	api := r.Group("/api/v1")
	{
		// 文件管理
		files := api.Group("/files")
		{
			files.POST("/upload", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "file uploaded"}) // 受多重限流
			})
			files.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "file data"}) // 只受 IP 限流
			})
			files.DELETE("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "file deleted"}) // 受 DELETE 限流
			})
		}

		// 图片管理
		images := api.Group("/images")
		{
			images.POST("/upload", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "image uploaded"}) // 受多重限流
			})
			images.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "image data"}) // 只受 IP 限流
			})
		}
	}

	r.Run(":8080")
}

// Example8_CustomErrorHandlers 自定义错误处理
func Example8_CustomErrorHandlers() {
	r := gin.New()

	r.Use(NewBuilder().
		WithIPLimit(10, time.Second).
		WithMethodLimit([]string{"POST"}, 1, time.Second).
		// 自定义 IP 限流错误
		WithIPErrorHandler(func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "请求过于频繁，请稍后再试",
				"retry_after": "1秒",
				"ip":          c.ClientIP(),
			})
			c.Abort()
		}).
		// 自定义方法限流错误
		WithMethodErrorHandler(func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "写操作过于频繁",
				"method":      c.Request.Method,
				"retry_after": "1秒",
			})
			c.Abort()
		}).
		WithCleanup(1 * time.Minute).
		Middleware())

	r.POST("/api/create", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "created"})
	})

	r.Run(":8080")
}

// Example9_LayeredLimits 分层限流（不同层级不同策略）
func Example9_LayeredLimits() {
	r := gin.New()

	// 1. 全局基础限流
	r.Use(NewBuilder().
		WithGlobalLimit(1000, time.Second).
		WithIPLimit(100, time.Second).
		WithCleanup(1 * time.Minute).
		Middleware())

	// 2. API 路由组的写操作限流
	apiLimiter := NewBuilder().
		WithMethodLimit([]string{"POST", "PUT", "DELETE"}, 10, time.Second).
		WithCleanup(1 * time.Minute)

	api := r.Group("/api/v1")
	api.Use(apiLimiter.Middleware())
	{
		// 3. 上传接口的严格限流
		uploadLimiter := NewBuilder().
			WithRouteLimit("/api/v1/files/upload", 1, time.Second, "POST").
			WithCleanup(1 * time.Minute)

		files := api.Group("/files")
		{
			files.POST("/upload",
				uploadLimiter.Middleware(), // 额外的限流中间件
				func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "file uploaded"})
				},
			)
			files.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "file data"})
			})
			files.DELETE("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
			})
		}
	}

	r.Run(":8080")
}

// Example10_RealWorldUsage 真实项目使用示例
func Example10_RealWorldUsage() {
	r := gin.New()

	// 中间件配置
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 限流配置（推荐的生产环境配置）
	limiter := NewBuilder().
		// 服务器总体限流：防止过载
		WithGlobalLimit(5000, time.Second).
		// IP 基础限流：防止单个用户滥用
		WithIPLimit(100, time.Second).
		// 写操作限流：保护数据库和存储
		WithMethodLimit([]string{"POST", "PUT", "DELETE"}, 20, time.Second).
		// 敏感操作限流（指定 HTTP 方法）
		WithRouteLimit("/api/v1/auth/login", 5, time.Minute, "POST").     // 登录：5次/分钟
		WithRouteLimit("/api/v1/auth/register", 3, time.Minute, "POST").  // 注册：3次/分钟
		WithRouteLimit("/api/v1/files/upload", 10, time.Minute, "POST").  // 上传：10次/分钟
		WithRouteLimit("/api/v1/images/upload", 10, time.Minute, "POST"). // 图片上传：10次/分钟
		WithRouteLimit("/api/v1/admin/*", 30, time.Minute, "*").          // 管理接口：30次/分钟（所有方法）
		// 自定义错误处理
		WithIPErrorHandler(func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
				"data":    nil,
			})
			c.Abort()
		}).
		WithMethodErrorHandler(func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "操作过于频繁，请稍后再试",
				"data":    nil,
			})
			c.Abort()
		}).
		WithRouteErrorHandler(func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "该操作过于频繁，请稍后再试",
				"data":    nil,
			})
			c.Abort()
		}).
		// 定期清理不活跃的 IP 记录
		WithCleanup(2 * time.Minute)

	// 应用限流中间件
	r.Use(limiter.Middleware())

	// 定义路由
	api := r.Group("/api/v1")
	{
		// 认证相关
		auth := api.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "login success"})
			})
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "register success"})
			})
		}

		// 文件管理
		files := api.Group("/files")
		{
			files.POST("/upload", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "file uploaded"})
			})
			files.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "file data"})
			})
			files.DELETE("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
			})
		}

		// 图片管理
		images := api.Group("/images")
		{
			images.POST("/upload", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "image uploaded"})
			})
			images.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "image data"})
			})
		}
	}

	r.Run(":8080")
}
