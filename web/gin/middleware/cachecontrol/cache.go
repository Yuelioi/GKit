package cachecontrol

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CacheBuilder struct {
	maxAge    time.Duration
	isPrivate bool
	noStore   bool
	immutable bool
}

// NewBuilder 初始化默认 builder
func NewBuilder() *CacheBuilder {
	return &CacheBuilder{}
}

// WithMaxAge 设置最大缓存时长
func (b *CacheBuilder) WithMaxAge(d time.Duration) *CacheBuilder {
	b.maxAge = d
	return b
}

// Private 表示响应是用户私有缓存（不共享）
func (b *CacheBuilder) Private() *CacheBuilder {
	b.isPrivate = true
	return b
}

// NoStore 禁止缓存
func (b *CacheBuilder) NoStore() *CacheBuilder {
	b.noStore = true
	return b
}

// Immutable 表示响应内容不会变化
func (b *CacheBuilder) Immutable() *CacheBuilder {
	b.immutable = true
	return b
}

// Build 生成 gin.HandlerFunc
func (b *CacheBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		var directives []string

		switch {
		case b.noStore:
			directives = append(directives, "no-store")
		default:
			if b.isPrivate {
				directives = append(directives, "private")
			} else {
				directives = append(directives, "public")
			}
			if b.maxAge > 0 {
				directives = append(directives, fmt.Sprintf("max-age=%d", int(b.maxAge.Seconds())))
			}
			if b.immutable {
				directives = append(directives, "immutable")
			}
		}

		c.Header("Cache-Control", strings.Join(directives, ", "))
		c.Next()
	}
}

// Default：默认 5 分钟公共缓存
func Default() gin.HandlerFunc {
	return NewBuilder().
		WithMaxAge(5 * time.Minute).
		Build()
}
