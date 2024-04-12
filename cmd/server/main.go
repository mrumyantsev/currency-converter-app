package main

import (
	"flag"

	"github.com/mrumyantsev/currency-converter-app/internal/app/server"
	"github.com/mrumyantsev/logx/log"
)

func main() {
	app, err := server.New()
	if err != nil {
		log.Fatal("failed to initialize application", err)
	}

	if isUserWantSave() {
		if err = app.SaveCurrencyDataToFile(); err != nil {
			log.Fatal("failed to save currencies to file", err)
		}

		return
	}

	if err = app.Run(); err != nil {
		log.Fatal("failed to run application", err)
	}
}

func isUserWantSave() bool {
	f := flag.Bool("s", false, "Save currency data to a local file")

	flag.Parse()

	return *f
}
