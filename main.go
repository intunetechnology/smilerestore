package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// intialize logger
	log.Logger = log.Output(zerolog.ConsoleWriter{
		NoColor: true,
		Out:     os.Stderr,
	})

	log.Info().Msg("smilerestore executing...")
	log.Info().Str("author", "Zach Snyder").Str("company", "intunetech").Str("license", "GPL v3.0").Msg("general info")
	time.Sleep(100)
	// fill path variable with execution directory
	execpath, err := os.Getwd()
	if err != nil {
		log.Error().Msg(err.Error())
	}

	// establish path flag and parse input, default to execpath, record as full path
	// establish which patient program is executing on based on ID
	pathStr := flag.String("path", execpath, "desired file path where program will execute")
	patientId := flag.String("patient", "0", "ID of patient program will perform corrective action on")
	flag.Parse()
	*pathStr, err = filepath.Abs(*pathStr) // convert to absolute path
	if err != nil {
		log.Error().Msg(err.Error())
	}

	log.Info().Str("path", *pathStr).Msg("using")
	log.Info().Str("patient", *patientId).Msg("using")

	// obtain slice of image directory and loop
	workingDir, err := ioutil.ReadDir(*pathStr)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	log.Log().Msg("-------------------------------------")
	// if no patient was specified dont run, explain why/how
	if *patientId != "0" {
		log.Info().Msg("checking for matching patient id")

		// loop working directory slice
		// if patient id not zero look for matching patient id

		for _, subdir := range workingDir {
			if strings.Split(subdir.Name(), "_")[1] == *patientId {
				// patient found do work
				log.Info().Str("directory", subdir.Name()).Msg("patient found")
				// just in case, check if original image directory exists or if its a directory
				needsAction := checkDirectory(subdir.Name(), filepath.Join(*pathStr, subdir.Name()))
				if needsAction {
					// perform corrective action
					recoverFile(filepath.Join(*pathStr, subdir.Name()))
				} else {
					log.Info().Str("directory", subdir.Name()).Msg("no original images directory found, no work to be done")
				}
			}
		}
	} else {
		log.Warn().Msg("patient id not specified, please use --patient flag to specify")
		log.Warn().Msg("example: .\\smilerestore.exe --path=..\\to_be_recovered\\. --patient=2315 ")
	}
	log.Log().Msg("-------------------------------------")
	log.Info().Msg("done")
}

func checkDirectory(name, path string) bool {
	// function checks if specified directory contains a subdirectory containing files needing recovery
	fp := filepath.Join(path, "OriginalImages.XVA")
	file, err := os.Stat(fp)
	if os.IsNotExist(err) {
		return false
	}
	// check if dir
	if !file.IsDir() {
		log.Fatal().Str("path", fp).Msg("Error: OriginalImages.XVA is a file not a directory")
	}
	log.Warn().Str("path", fp).Msg("Recovery directory exists")

	return true
}

func compareFile(a, b string) bool {
	// function checks if two files a & b are the same data based on name
	// example: file.auto vs file_Original.auto
	c := strings.Split(a, ".")
	d := strings.Fields(b)

	if strings.ToLower(c[1]) == "auto" || strings.ToLower(c[1]) == "autoxvtag" {
		// log.Log().Msg(fmt.Sprintf("Comparing %s <-> %s", a, b))
		if c[0] == d[0] && c[1] == strings.Split(d[1], ".")[1] {
			return true
		}
	}
	return false
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func recoverFile(path string) string {
	// break down directory composition into a slice
	patientDir, err := ioutil.ReadDir(path)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	// break down OriginalImages.XVA into a slice
	recoveryDir, err := ioutil.ReadDir(filepath.Join(path, "OriginalImages.XVA"))
	if err != nil {
		log.Error().Msg(err.Error())
	}

	// comparison loop
	for _, a := range patientDir {
		for _, b := range recoveryDir {
			// log.Log().Msg(fmt.Sprintf("Comparing %s <-> %s", a.Name(), b.Name()))
			if compareFile(a.Name(), b.Name()) {
				// rename and move file
				log.Info().Str("corrupted", a.Name()).Str("original", b.Name()).Msg("match found, Attempting recovery")
				// store paths
				corrupt := filepath.Join(path, a.Name())
				original := filepath.Join(filepath.Join(path, "OriginalImages.XVA"), b.Name())
				log.Log().Str("path", corrupt).Msg("corrupt")
				log.Log().Str("path", original).Msg("original")
				// copy files
				bytes, err := copy(original, corrupt)
				if err != nil {
					log.Error().Str("event", "copy").Msg(err.Error())
				}
				log.Info().Str("event", "copy").Msg(fmt.Sprintf("success, copied %d bytes", bytes))
			}
		}
	}
	return "0"
}
