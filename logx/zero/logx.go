package zero

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type LogConfig struct {
	Level      zerolog.Level
	NoColor    bool
	WithCaller bool
	Output     *os.File
}

func New(cfg LogConfig) zerolog.Logger {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	cw := zerolog.ConsoleWriter{
		Out:        cfg.Output,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    cfg.NoColor,
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
	cw.FormatFieldName = func(i interface{}) string { return fmt.Sprintf("%s=", i) }
	cw.FormatFieldValue = func(i interface{}) string { return fmt.Sprintf("%s", i) }

	logger := zerolog.New(cw).With().Timestamp().Logger()
	if cfg.WithCaller {
		logger = logger.With().Caller().Logger()
	}
	zerolog.SetGlobalLevel(cfg.Level)

	return logger
}
