package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("smilerestore executing...")
	log.Info().Str("author", "Zach Snyder").Str("company", "intunetech").Str("license", "GPL v3.0").Msg("general info")

	path, err := os.Getwd()

	if err != nil {
		log.Error().Msg(err.Error())
	}

	pathStr := flag.String("path", path, "desired file path where program will execute")

	log.Info().Str("path", *pathStr).Msg("execution path ")

}

func checkDirectory(dir string) (string, error) {
	return "0", nil
}

func recoverFile(filename string) (string, error) {
	return "0", nil
}
