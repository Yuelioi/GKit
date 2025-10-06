package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 支持 IP 和全局限流
type Builder struct {
	// IP 限流
	ipEnabled  bool
	ipMax      int
	ipInterval time.Duration

	// 全局限流
	globalEnabled  bool
	globalMax      int
	globalInterval time.Duration

	cleanupInterval time.Duration

	// 自定义限流错误处理
	IPErrorHandler     func(c *gin.Context)
	GlobalErrorHandler func(c *gin.Context)

	// 内部状态
	clients map[string]*client
	mu      sync.Mutex
	global  *rate.Limiter
}

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func New() *Builder {
	return &Builder{
		clients: make(map[string]*client),
		// 默认错误处理
		IPErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "IP rate limit exceeded"})
			c.Abort()
		},
		GlobalErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Global rate limit exceeded"})
			c.Abort()
		},
	}
}

func (rl *Builder) WithIPLimit(max int, interval time.Duration) *Builder {
	if max > 0 {
		rl.ipEnabled = true
		rl.ipMax = max
		rl.ipInterval = interval
	}
	return rl
}

func (rl *Builder) WithGlobalLimit(max int, interval time.Duration) *Builder {
	if max > 0 {
		rl.globalEnabled = true
		rl.globalMax = max
		rl.globalInterval = interval
		rl.global = rate.NewLimiter(rate.Every(interval/time.Duration(max)), max)
	}
	return rl
}

// 设置自定义 IP 错误处理
func (rl *Builder) WithIPErrorHandler(fn func(c *gin.Context)) *Builder {
	if fn != nil {
		rl.IPErrorHandler = fn
	}
	return rl
}

// 设置自定义全局错误处理
func (rl *Builder) WithGlobalErrorHandler(fn func(c *gin.Context)) *Builder {
	if fn != nil {
		rl.GlobalErrorHandler = fn
	}
	return rl
}

func (rl *Builder) WithCleanup(interval time.Duration) *Builder {
	if interval > 0 {
		rl.cleanupInterval = interval
	}
	return rl
}

// Middleware 返回 gin.HandlerFunc
func (rl *Builder) Middleware() gin.HandlerFunc {
	if rl.ipEnabled && rl.cleanupInterval > 0 {
		go rl.cleanupClients()
	}

	return func(c *gin.Context) {
		// 全局限流
		if rl.globalEnabled && rl.global != nil {
			if !rl.global.Allow() {
				rl.GlobalErrorHandler(c)
				return
			}
		}

		// IP 限流
		if rl.ipEnabled {
			ip := c.ClientIP()
			rl.mu.Lock()
			cl, exists := rl.clients[ip]
			if !exists {
				cl = &client{
					limiter: rate.NewLimiter(rate.Every(rl.ipInterval/time.Duration(rl.ipMax)), rl.ipMax),
				}
				rl.clients[ip] = cl
			}
			cl.lastSeen = time.Now()
			rl.mu.Unlock()

			if !cl.limiter.Allow() {
				rl.IPErrorHandler(c)
				return
			}
		}

		c.Next()
	}
}

// 清理长时间未访问的 IP
func (rl *Builder) cleanupClients() {
	interval := rl.cleanupInterval
	if interval <= 0 {
		interval = time.Minute
	}
	for {
		time.Sleep(interval)
		rl.mu.Lock()
		for ip, cl := range rl.clients {
			if time.Since(cl.lastSeen) > 5*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}
