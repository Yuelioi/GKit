package server

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ServerConfig 配置服务器参数
type ServerConfig struct {
	Addr        string
	Logger      zerolog.Logger
	Mode        string // release 或 debug
	APIPrefix   string
	Middlewares []gin.HandlerFunc
	SPAPath     string
	IgnorePaths []string
	EnableCORS  bool
}

// Start 启动服务器
func Start(cfg ServerConfig, registerRoutes func(api *gin.RouterGroup)) error {
	if cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		// 开发模式，支持跨域
		if cfg.EnableCORS {
			cfg.Middlewares = append(cfg.Middlewares, cors.Default())
		}
		// 默认监听 0.0.0.0
		if len(cfg.Addr) > 0 && cfg.Addr[0] == ':' {
			cfg.Addr = "0.0.0.0" + cfg.Addr
		}
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(cfg.Middlewares...)

	// API 路由组
	api := r.Group(cfg.APIPrefix)
	registerRoutes(api)

	// 静态资源 + SPA
	if cfg.SPAPath != "" {
		// 静态资源中间件，只处理非 API 请求
		r.Use(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, cfg.APIPrefix) {
				c.Next()
				return
			}
			// 尝试返回静态文件
			static.Serve("/", static.LocalFile(cfg.SPAPath, false))(c)
			// 如果文件存在，static 会处理完，不再继续
			if c.IsAborted() {
				return
			}
			// 如果文件不存在，继续执行 NoRoute fallback
			c.Next()
		})

		// SPA fallback，前端刷新路由
		r.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, cfg.APIPrefix) {
				c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
				return
			}
			c.File(cfg.SPAPath + "/index.html")
		})
	}

	cfg.Logger.Info().Str("addr", cfg.Addr).Msg("服务启动")
	return r.Run(cfg.Addr)
}
