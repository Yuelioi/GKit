package server

import (
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
	} else {
		// 开发模式，支持跨域
		if cfg.EnableCORS {
			cfg.Middlewares = append(cfg.Middlewares, cors.Default())
		}
		// 默认监听 0.0.0.0
		if len(cfg.Addr) > 0 && cfg.Addr[0] == ':' {
			cfg.Addr = "0.0.0.0" + cfg.Addr
		}
	}

	r := gin.New()
	r.Use(cfg.Middlewares...)

	// API 路由组
	api := r.Group(cfg.APIPrefix)
	registerRoutes(api)

	// 静态资源 + SPA
	if cfg.SPAPath != "" {
		r.Use(static.Serve("/", static.LocalFile(cfg.SPAPath, true)))
	}

	cfg.Logger.Info().Str("addr", cfg.Addr).Msg("服务启动")
	return r.Run(cfg.Addr)
}
