package zerologx

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type LoggerBuilder struct {
	level      zerolog.Level
	withCaller bool
	noColor    bool
	output     *os.File
}

func Default() zerolog.Logger {
	return NewBuilder().
		Level(zerolog.DebugLevel). // 默认 Debug
		Output(os.Stdout).         // 默认终端
		Build()
}

// 初始化 builder
func NewBuilder() *LoggerBuilder {
	return &LoggerBuilder{
		level:  zerolog.InfoLevel,
		output: os.Stdout,
	}
}

func (b *LoggerBuilder) Level(l zerolog.Level) *LoggerBuilder {
	b.level = l
	return b
}

func (b *LoggerBuilder) WithCaller() *LoggerBuilder {
	b.withCaller = true
	return b
}

func (b *LoggerBuilder) NoColor() *LoggerBuilder {
	b.noColor = true
	return b
}

func (b *LoggerBuilder) Output(f *os.File) *LoggerBuilder {
	b.output = f
	return b
}

func (b *LoggerBuilder) Build() zerolog.Logger {
	var logger zerolog.Logger
	isConsole := b.output == os.Stdout || b.output == os.Stderr

	if isConsole {
		// 终端输出，带颜色
		cw := zerolog.ConsoleWriter{
			Out:        b.output,
			TimeFormat: "2006-01-02 15:04:05",
			NoColor:    b.noColor,
		}

		colors := map[string]string{
			"debug": "\x1b[35mDEBG\x1b[0m",
			"info":  "\x1b[32mINFO\x1b[0m",
			"warn":  "\x1b[33mWARN\x1b[0m",
			"error": "\x1b[31mERRO\x1b[0m",
			"fatal": "\x1b[31mFATL\x1b[0m",
			"panic": "\x1b[31mPANC\x1b[0m",
		}

		cw.FormatLevel = func(i interface{}) string {
			if ll, ok := i.(string); ok {
				if color, found := colors[ll]; found {
					return fmt.Sprintf("[%s]", color)
				}
				return fmt.Sprintf("[%s]", strings.ToUpper(ll))
			}
			return fmt.Sprintf("[%s]", i)
		}

		// 直接使用 ConsoleWriter，不要再包装
		logger = zerolog.New(cw).With().Timestamp().Logger()
	} else {
		// 文件或其他 io.Writer，使用 JSON 格式
		logger = zerolog.New(b.output).With().Timestamp().Logger()
	}

	if b.withCaller {
		logger = logger.With().Caller().Logger()
	}

	logger = logger.Level(b.level)
	return logger
}
