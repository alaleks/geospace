package main

import (
	"github.com/alaleks/geospace/internal/client"
)

// application information
var (
	Version string
	Host    string
	Name    string
)

func main() {
	cl := client.New(Version, Host, Name)
	cl.Run()
}
