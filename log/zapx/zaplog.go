package zapx

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerBuilder struct {
	level      zapcore.Level
	withCaller bool
	noColor    bool
	output     *os.File
}

func Default() *zap.Logger {
	return NewBuilder().
		Level(zapcore.DebugLevel).
		Output(os.Stdout).
		Build()
}

// 初始化 builder
func NewBuilder() *LoggerBuilder {
	return &LoggerBuilder{
		level:  zapcore.InfoLevel,
		output: os.Stdout,
	}
}

func (b *LoggerBuilder) Level(l zapcore.Level) *LoggerBuilder {
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

func (b *LoggerBuilder) Build() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    b.levelEncoder(),
		EncodeTime:     zapcore.TimeEncoderOfLayout("01-02 15:04:05"),
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(b.output),
		b.level,
	)

	logger := zap.New(core)

	if b.withCaller {
		logger = logger.WithOptions(zap.AddCaller())
	}

	return logger
}

func (b *LoggerBuilder) levelEncoder() zapcore.LevelEncoder {
	if b.noColor {
		return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(fmt.Sprintf("[%s]", l.CapitalString()))
		}
	}

	colors := map[zapcore.Level]string{
		zapcore.DebugLevel:  "\x1b[35m", // 紫色
		zapcore.InfoLevel:   "\x1b[32m", // 绿色
		zapcore.WarnLevel:   "\x1b[33m", // 黄色
		zapcore.ErrorLevel:  "\x1b[31m", // 红色
		zapcore.DPanicLevel: "\x1b[31m", // 红色
		zapcore.PanicLevel:  "\x1b[31m", // 红色
		zapcore.FatalLevel:  "\x1b[31m", // 红色
	}

	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		color := colors[l]
		reset := "\x1b[0m"
		enc.AppendString(fmt.Sprintf("[%s%s%s]", color, l.CapitalString(), reset))
	}
}
