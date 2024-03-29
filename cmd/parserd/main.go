package main

import (
	"flag"

	"github.com/mrumyantsev/currency-converter-app/internal/app/parserd"
)

func main() {
	app := parserd.New()

	if isUserWantSave() {
		app.SaveCurrencyDataToFile()

		return
	}

	app.Run()
}

func isUserWantSave() bool {
	f := flag.Bool("s", false, "Save currency data to a local file")

	flag.Parse()

	return *f
}
