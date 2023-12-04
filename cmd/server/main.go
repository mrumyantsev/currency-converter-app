package main

import (
	"os"

	"github.com/mrumyantsev/currency-converter/internal/app/server"
)

func main() {
	server := server.New()

	if isGotSecondArg("--save") {
		server.SaveCurrencyDataToFile()
	} else {
		server.Run()
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
