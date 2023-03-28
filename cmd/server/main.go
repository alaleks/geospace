package main

import "github.com/alaleks/geospace/internal/server/app"

func main() {
	server := app.New()
	server.Run()
}
