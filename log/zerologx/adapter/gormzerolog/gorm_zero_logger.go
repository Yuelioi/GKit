package gormx

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

// GORM Zerolog 日志实现
type Logger struct {
	log           zerolog.Logger
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

// WithZerolog 设置 zerolog 实例
func WithZerolog(zl zerolog.Logger) Option {
	return func(l *Logger) {
		l.log = zl
	}
}

// 创建一个 GORM Zerolog Logger
func New(opts ...Option) *Logger {
	// 默认值
	l := &Logger{
		log:           zerolog.New(nil),       // 默认输出 nil（无输出）
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
		l.log.Info().Msgf(msg, data...)
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Warn {
		l.log.Warn().Msgf(msg, data...)
	}
}

func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Error {
		l.log.Error().Msgf(msg, data...)
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level == logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	event := l.log.With().
		Str("sql", sql).
		Dur("elapsed", elapsed).
		Int64("rows", rows).
		Logger()

	switch {
	case err != nil:
		event.Error().Err(err).Msg("gorm error")
	case elapsed > l.slowThreshold:
		event.Warn().Msg("slow query")
	default:
		if l.level >= logger.Info {
			event.Debug().Msg("gorm query")
		}
	}
}
