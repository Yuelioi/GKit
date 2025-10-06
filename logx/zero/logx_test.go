package zero

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestLoggerNew(t *testing.T) {

	log := New(LogConfig{Level: zerolog.DebugLevel, WithCaller: false})
	log.Info().Str("module", "core").Msg("Logger initialized")
	log.Warn().Msg("This is a warning")

	log2 := New(LogConfig{Level: zerolog.DebugLevel, WithCaller: false})
	log2.Info().Str("module", "core").Msg("Logger initialized")

	// log.Panic().Msg("This is a panic")

	// Output:
	// 2025-10-06 18:53:13 [INFO] example_test.go:8 > Logger initialized module=core
	// 2025-10-06 18:53:13 [WARN] example_test.go:9 > This is a warning

	t.Log("✅ Logger initialization")

}
