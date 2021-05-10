package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// intialize logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("smilerestore executing...")
	log.Info().Str("author", "Zach Snyder").Str("company", "intunetech").Str("license", "GPL v3.0").Msg("general info")

	// fill path variable with execution directory
	execpath, err := os.Getwd()
	if err != nil {
		log.Error().Msg(err.Error())
	}

	// establish path flag and parse input, default to execpath, record as full path
	pathStr := flag.String("path", execpath, "desired file path where program will execute")
	flag.Parse()
	*pathStr, err = filepath.Abs(*pathStr) // convert to absolute path
	if err != nil {
		log.Error().Msg(err.Error())
	}

	log.Info().Str("path", *pathStr).Msg("execution path")

	// obtain slice of image directory and loop
	workingDir, err := ioutil.ReadDir(*pathStr)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	// loop working directory slice
	for _, subdir := range workingDir {
		log.Info().Str("name", subdir.Name()).Msg("checking dir")
		checkDirectory(subdir.Name(), filepath.Join(*pathStr, subdir.Name()))
	}

}

func checkDirectory(name string, path string) (bool, error) {
	// function checks if specified directory contains a subdirectory containing files needing recovery
	fp := filepath.Join(path, "OriginalImages.XVA")
	file, err := os.Stat(fp)
	if os.IsNotExist(err) {
		return false, nil
	}
	log.Warn().Str("path", fp).Msg("Recovery directory exists")
	// check if dir
	if !file.IsDir() {
		log.Error().Str("path", fp).Msg("Error: OriginalImages.XVA is a file not a directory")
	}

	// report back
	// log.Log().Msg("test")
	return false, nil
}

func recoverFile(filename string) (string, error) {
	return "0", nil

}
