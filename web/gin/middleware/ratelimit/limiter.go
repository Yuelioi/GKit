package ratelimit

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 支持 IP、全局、路由、方法级别的限流
type Builder struct {
	// IP 限流
	ipEnabled  bool
	ipMax      int
	ipInterval time.Duration

	// 全局限流
	globalEnabled  bool
	globalMax      int
	globalInterval time.Duration

	// 路由级别限流规则
	routeRules map[string]*RouteRule

	// 方法级别限流规则（如只限制 POST/DELETE）
	methodRules map[string]*MethodRule

	// 全局路由限流（所有 IP 共享）
	globalRouteRules map[string]*rate.Limiter

	cleanupInterval time.Duration

	// 自定义限流错误处理
	IPErrorHandler     func(c *gin.Context)
	GlobalErrorHandler func(c *gin.Context)
	RouteErrorHandler  func(c *gin.Context)
	MethodErrorHandler func(c *gin.Context)

	// 内部状态
	clients       map[string]*client            // IP 限流
	routeClients  map[string]map[string]*client // 路由限流: route -> ip -> client
	methodClients map[string]map[string]*client // 方法限流: method -> ip -> client
	mu            sync.Mutex
	global        *rate.Limiter
	cleanupOnce   sync.Once // 确保清理 goroutine 只启动一次
}

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RouteRule 路由级别的限流规则
type RouteRule struct {
	Max      int           // 每个 IP 在该路由下的最大请求数
	Interval time.Duration // 时间窗口
	Pattern  string        // 路由匹配模式（支持前缀匹配）
	Methods  []string      // 限制的 HTTP 方法，空或包含 "*" 表示所有方法
}

// MethodRule 方法级别的限流规则
type MethodRule struct {
	Max      int           // 每个 IP 对该方法的最大请求数
	Interval time.Duration // 时间窗口
	Methods  []string      // 限制的方法列表
}

// 默认限流：全局 100 次/秒，每个 IP 10 次/秒
func Default() gin.HandlerFunc {
	return NewBuilder().
		WithGlobalLimit(100, time.Second).
		WithIPLimit(10, time.Second).
		WithCleanup(1 * time.Minute).
		Middleware()
}

func NewBuilder() *Builder {
	return &Builder{
		clients:          make(map[string]*client),
		routeClients:     make(map[string]map[string]*client),
		methodClients:    make(map[string]map[string]*client),
		routeRules:       make(map[string]*RouteRule),
		methodRules:      make(map[string]*MethodRule),
		globalRouteRules: make(map[string]*rate.Limiter),
		// 默认错误处理
		IPErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "IP rate limit exceeded",
				"retry_after": "1s",
			})
			c.Abort()
		},
		GlobalErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Global rate limit exceeded",
				"retry_after": "1s",
			})
			c.Abort()
		},
		RouteErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Route rate limit exceeded",
				"path":        c.Request.URL.Path,
				"retry_after": "1s",
			})
			c.Abort()
		},
		MethodErrorHandler: func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Method rate limit exceeded",
				"method":      c.Request.Method,
				"retry_after": "1s",
			})
			c.Abort()
		},
	}
}

// WithIPLimit 设置全局 IP 限流
func (rl *Builder) WithIPLimit(max int, interval time.Duration) *Builder {
	if max > 0 {
		rl.ipEnabled = true
		rl.ipMax = max
		rl.ipInterval = interval
	}
	return rl
}

// WithGlobalLimit 设置全局限流
func (rl *Builder) WithGlobalLimit(max int, interval time.Duration) *Builder {
	if max > 0 {
		rl.globalEnabled = true
		rl.globalMax = max
		rl.globalInterval = interval
		rl.global = rate.NewLimiter(rate.Every(interval/time.Duration(max)), max)
	}
	return rl
}

// WithRouteLimit 为特定路由设置限流
// pattern 支持前缀匹配，如 "/api/v1/files/upload" 或 "/api/v1/images/*"
// methods 指定限制的 HTTP 方法，如 []string{"POST"}，空或 []string{"*"} 表示所有方法
func (rl *Builder) WithRouteLimit(pattern string, max int, interval time.Duration, methods ...string) *Builder {
	if max > 0 {
		// 如果没有指定方法或者包含 "*"，则表示所有方法
		if len(methods) == 0 {
			methods = []string{"*"}
		}

		rl.routeRules[pattern] = &RouteRule{
			Max:      max,
			Interval: interval,
			Pattern:  pattern,
			Methods:  methods,
		}
	}
	return rl
}

// WithMethodLimit 为特定 HTTP 方法设置限流
// methods: 如 "POST", "DELETE", "PUT" 等
func (rl *Builder) WithMethodLimit(methods []string, max int, interval time.Duration) *Builder {
	if max > 0 && len(methods) > 0 {
		key := strings.Join(methods, ",")
		rl.methodRules[key] = &MethodRule{
			Max:      max,
			Interval: interval,
			Methods:  methods,
		}
	}
	return rl
}

// WithIPErrorHandler 设置自定义 IP 错误处理
func (rl *Builder) WithIPErrorHandler(fn func(c *gin.Context)) *Builder {
	if fn != nil {
		rl.IPErrorHandler = fn
	}
	return rl
}

// WithGlobalErrorHandler 设置自定义全局错误处理
func (rl *Builder) WithGlobalErrorHandler(fn func(c *gin.Context)) *Builder {
	if fn != nil {
		rl.GlobalErrorHandler = fn
	}
	return rl
}

// WithRouteErrorHandler 设置自定义路由错误处理
func (rl *Builder) WithRouteErrorHandler(fn func(c *gin.Context)) *Builder {
	if fn != nil {
		rl.RouteErrorHandler = fn
	}
	return rl
}

// WithMethodErrorHandler 设置自定义方法错误处理
func (rl *Builder) WithMethodErrorHandler(fn func(c *gin.Context)) *Builder {
	if fn != nil {
		rl.MethodErrorHandler = fn
	}
	return rl
}

// WithCleanup 设置清理间隔
func (rl *Builder) WithCleanup(interval time.Duration) *Builder {
	if interval > 0 {
		rl.cleanupInterval = interval
	}
	return rl
}

// Middleware 返回 gin.HandlerFunc
func (rl *Builder) Middleware() gin.HandlerFunc {
	// 确保清理 goroutine 只启动一次
	if (rl.ipEnabled || len(rl.routeRules) > 0 || len(rl.methodRules) > 0) && rl.cleanupInterval > 0 {
		rl.startCleanupOnce()
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 1. 全局限流
		if rl.globalEnabled && rl.global != nil {
			if !rl.global.Allow() {
				rl.GlobalErrorHandler(c)
				return
			}
		}

		// 2. 方法级别限流（如只限制 POST/DELETE）
		if len(rl.methodRules) > 0 {
			for _, rule := range rl.methodRules {
				if rl.matchMethod(method, rule.Methods) {
					if !rl.checkMethodLimit(ip, method, rule) {
						rl.MethodErrorHandler(c)
						return
					}
				}
			}
		}

		// 3. 路由级别限流
		if len(rl.routeRules) > 0 {
			for pattern, rule := range rl.routeRules {
				if rl.matchRoute(path, pattern) && rl.matchRouteMethod(method, rule.Methods) {
					if !rl.checkRouteLimit(ip, pattern, rule) {
						rl.RouteErrorHandler(c)
						return
					}
				}
			}
		}

		// 4. IP 全局限流
		if rl.ipEnabled {
			if !rl.checkIPLimit(ip) {
				rl.IPErrorHandler(c)
				return
			}
		}

		c.Next()
	}
}

// checkIPLimit 检查 IP 限流
func (rl *Builder) checkIPLimit(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cl, exists := rl.clients[ip]
	if !exists {
		cl = &client{
			limiter: rate.NewLimiter(rate.Every(rl.ipInterval/time.Duration(rl.ipMax)), rl.ipMax),
		}
		rl.clients[ip] = cl
	}
	cl.lastSeen = time.Now()

	return cl.limiter.Allow()
}

// checkRouteLimit 检查路由限流
func (rl *Builder) checkRouteLimit(ip, pattern string, rule *RouteRule) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.routeClients[pattern] == nil {
		rl.routeClients[pattern] = make(map[string]*client)
	}

	cl, exists := rl.routeClients[pattern][ip]
	if !exists {
		cl = &client{
			limiter: rate.NewLimiter(rate.Every(rule.Interval/time.Duration(rule.Max)), rule.Max),
		}
		rl.routeClients[pattern][ip] = cl
	}
	cl.lastSeen = time.Now()

	return cl.limiter.Allow()
}

// checkMethodLimit 检查方法限流
func (rl *Builder) checkMethodLimit(ip, method string, rule *MethodRule) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.methodClients[method] == nil {
		rl.methodClients[method] = make(map[string]*client)
	}

	cl, exists := rl.methodClients[method][ip]
	if !exists {
		cl = &client{
			limiter: rate.NewLimiter(rate.Every(rule.Interval/time.Duration(rule.Max)), rule.Max),
		}
		rl.methodClients[method][ip] = cl
	}
	cl.lastSeen = time.Now()

	return cl.limiter.Allow()
}

// matchRoute 匹配路由（支持前缀匹配）
func (rl *Builder) matchRoute(path, pattern string) bool {
	// 精确匹配
	if path == pattern {
		return true
	}
	// 前缀匹配
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(path, prefix)
	}
	return strings.HasPrefix(path, pattern)
}

// matchMethod 匹配 HTTP 方法
func (rl *Builder) matchMethod(method string, methods []string) bool {
	for _, m := range methods {
		if m == "*" || strings.EqualFold(method, m) {
			return true
		}
	}
	return false
}

// matchRouteMethod 匹配路由规则的 HTTP 方法
func (rl *Builder) matchRouteMethod(method string, methods []string) bool {
	// 如果 methods 为空或包含 "*"，匹配所有方法
	if len(methods) == 0 {
		return true
	}
	for _, m := range methods {
		if m == "*" || strings.EqualFold(method, m) {
			return true
		}
	}
	return false
}

// startCleanupOnce 确保清理 goroutine 只启动一次
func (rl *Builder) startCleanupOnce() {
	rl.cleanupOnce.Do(func() {
		go rl.cleanupClients()
	})
}

// cleanupClients 清理长时间未访问的 IP
func (rl *Builder) cleanupClients() {
	interval := rl.cleanupInterval
	if interval <= 0 {
		interval = time.Minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()

		// 清理 IP 限流
		for ip, cl := range rl.clients {
			if time.Since(cl.lastSeen) > 5*time.Minute {
				delete(rl.clients, ip)
			}
		}

		// 清理路由限流
		for route, clients := range rl.routeClients {
			for ip, cl := range clients {
				if time.Since(cl.lastSeen) > 5*time.Minute {
					delete(clients, ip)
				}
			}
			if len(clients) == 0 {
				delete(rl.routeClients, route)
			}
		}

		// 清理方法限流
		for method, clients := range rl.methodClients {
			for ip, cl := range clients {
				if time.Since(cl.lastSeen) > 5*time.Minute {
					delete(clients, ip)
				}
			}
			if len(clients) == 0 {
				delete(rl.methodClients, method)
			}
		}

		rl.mu.Unlock()
	}
}
