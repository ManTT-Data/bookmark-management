package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// SetLogLevel sets the global log level based on the environment variable LOG_LEVEL
func SetLogLevel() {
	level, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil || level == zerolog.NoLevel {
		level = zerolog.InfoLevel // fallback
	}
	zerolog.SetGlobalLevel(level)
}
