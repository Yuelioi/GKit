package zerologx

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func TestLoggerNew(t *testing.T) {
	logger := NewBuilder().
		Level(zerolog.DebugLevel).
		WithCaller().
		Output(os.Stdout).
		Build()

	logger.Info().Str("module", "core").Msg("Logger initialized")
	logger.Warn().Msg("This is a warning")

	f, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	logger = NewBuilder().
		Level(zerolog.InfoLevel).
		WithCaller().
		Output(f).
		Build()

	logger.Info().
		Str("module", "core").
		Msg("Logger initialized")

	logger.Warn().
		Str("module", "core").
		Msg("This is a warning")

	t.Log("✅ Logger initialization")

}
