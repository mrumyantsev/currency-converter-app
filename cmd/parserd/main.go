package main

import (
	"os"

	"github.com/mrumyantsev/currency-converter/internal/app/parserd"
)

func main() {
	parserD := parserd.New()

	if isGotSecondArg("--save") {
		parserD.SaveCurrencyDataToFile()
	} else {
		parserD.Run()
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
