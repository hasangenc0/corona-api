package main

import (
	"github.com/hasangenc0/corona/pkg/server"
)

func main() {
	app := &server.Server{}
	app.Bootstrap().Start()
}
