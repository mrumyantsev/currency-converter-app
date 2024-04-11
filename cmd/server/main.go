package main

import "github.com/mrumyantsev/currency-converter-app/internal/app/server"

func main() {
	app := server.New()

	app.Run()
}
