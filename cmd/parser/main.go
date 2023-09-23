package main

import (
	"os"

	"github.com/mrumyantsev/currency-converter/internal/app/parser"
)

func main() {
	parser := parser.New()

	// // } if isArgFound("--save") {
	// // 	parser.SaveCurrencyDataFile()
	// // } else {
	// // 	parser.Run()
	// // }

	parser.DoWithDB()
}

func isArgFound(arg string) bool {
	for _, inputArg := range os.Args {
		if inputArg == arg {
			return true
		}
	}

	return false
}
