package main

import (
	"flag"
	"os"
	"time"

	"github.com/mrumyantsev/currency-converter-app/internal/app/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	conWrt := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	log.Logger = log.Output(conWrt)
}

func main() {
	app, err := server.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize application")
	}

	if isUserWantSave() {
		if err = app.SaveCurrencyDataToFile(); err != nil {
			log.Fatal().Err(err).Msg("failed to save currencies to file")
		}

		return
	}

	if err = app.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run application")
	}
}

func isUserWantSave() bool {
	f := flag.Bool("s", false, "Save currency data to a local file")

	flag.Parse()

	return *f
}
