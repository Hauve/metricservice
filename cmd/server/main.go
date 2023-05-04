package main

import (
	"fmt"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/logger"
	"github.com/Hauve/metricservice.git/internal/server"
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go catchCloseSignals()

	cfg := config.LoadServerConfig()
	st := storage.NewMemStorage()
	r := chi.NewRouter()
	lg, err := logger.New()
	if err != nil {
		log.Fatalf("Logger creating failed: %s", err)
	}
	serv := server.New(cfg, st, r, lg)
	serv.Run()
}

// Для завершения процесса по ctrl+z в linux терминале
func catchCloseSignals() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGTSTP)
	go func() {
		sig := <-c
		fmt.Printf("\nGot %s signal. Aborting...\n", sig)
		os.Exit(1)
	}()
}
