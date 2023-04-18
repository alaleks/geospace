package main

import (
	"github.com/alaleks/geospace/internal/cli"
)

// application information
var (
	Version string
	Host    string
	Name    string
	Token   string
)

func main() {
	cli.UpdateAppInfo(Version, Host, Name, Token)
	cli.Execute()
}
