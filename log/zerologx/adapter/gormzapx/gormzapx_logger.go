package gormzapx

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// GORM Zap 日志实现
type Logger struct {
	log           *zap.Logger
	level         logger.LogLevel
	slowThreshold time.Duration
}

// Option 定义配置函数类型
type Option func(*Logger)

// WithLevel 设置日志级别
func WithLevel(level logger.LogLevel) Option {
	return func(l *Logger) {
		l.level = level
	}
}

// WithSlowThreshold 设置慢查询阈值
func WithSlowThreshold(threshold time.Duration) Option {
	return func(l *Logger) {
		l.slowThreshold = threshold
	}
}

// WithZap 设置 zap 实例
func WithZap(zl *zap.Logger) Option {
	return func(l *Logger) {
		l.log = zl
	}
}

// 创建一个 GORM Zap Logger
func New(opts ...Option) *Logger {
	// 默认值
	l := &Logger{
		log:           zap.NewNop(),           // 默认 noop logger
		level:         logger.Info,            // 默认 Info 级别
		slowThreshold: 200 * time.Millisecond, // 默认 200ms
	}

	// 应用 Option
	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	l.level = level
	return l
}

func (l *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		if len(data) == 0 {
			l.log.Info(msg)
		} else {
			l.log.Info(msg, zap.Any("data", data))
		}
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Warn {
		if len(data) == 0 {
			l.log.Warn(msg)
		} else {
			l.log.Warn(msg, zap.Any("data", data))
		}
	}
}

func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Error {
		if len(data) == 0 {
			l.log.Error(msg)
		} else {
			l.log.Error(msg, zap.Any("data", data))
		}
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level == logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	}

	switch {
	case err != nil:
		l.log.Error("gorm error", append(fields, zap.Error(err))...)
	case elapsed > l.slowThreshold:
		l.log.Warn("slow query", fields...)
	default:
		if l.level >= logger.Info {
			l.log.Debug("gorm query", fields...)
		}
	}
}
