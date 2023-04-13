package main

import (
	"github.com/Hauve/metricservice.git/internal/handlers"
	"github.com/Hauve/metricservice.git/internal/server"
)

func main() {
	serv := server.New(*handlers.New())
	serv.Run()
}
