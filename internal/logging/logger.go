package logging

import (
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2024-01-01 15:04:05"}
	Log = zerolog.New(output).With().Timestamp().Logger()

	// Set global log level (can be adjusted based on environment or config)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
