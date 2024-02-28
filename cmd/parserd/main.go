package main

import (
	"os"

	"github.com/mrumyantsev/currency-converter-app/internal/app/parserd"
)

func main() {
	app := parserd.New()

	if isGotSecondArg("--save") {
		app.SaveCurrencyDataToFile()
	} else {
		app.Run()
	}
}

func isGotSecondArg(arg string) bool {
	if len(os.Args) != 2 {
		return false
	}

	if os.Args[1] == arg {
		return true
	}

	return false
}
