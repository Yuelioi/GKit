package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ServerConfig 可配置启动参数
type ServerConfig struct {
	Addr        string
	Logger      zerolog.Logger
	Mode        string
	APIPrefix   string
	Middlewares []gin.HandlerFunc
}

// Start 启动服务器
func Start(cfg ServerConfig, setupRoutes func(r *gin.RouterGroup)) error {
	if cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(cfg.Middlewares...)
	r.Use(gin.Recovery())

	api := r.Group(cfg.APIPrefix)
	setupRoutes(api)

	return r.Run(cfg.Addr)
}
