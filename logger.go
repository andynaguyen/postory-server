package postory_server

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	zerolog.Logger
}

func NewLogger() Logger {
	if os.Getenv("STAGE") == "prod" {
		return Logger{log.Logger}
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		return Logger{log.Logger}
	}
}
