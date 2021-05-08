package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Str("foo", "bar").Msg("Hello world")
}

// Output: {"time":1516134303,"level":"debug","message":"hello world"}
