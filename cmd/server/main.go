package main

import (
	"net/http"

	"github.com/Hauve/metricservice.git/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	service := handlers.New()

	r := chi.NewRouter()

	r.Get("/value/{metricType}/{metricName}", service.GetHandler)
	r.Get("/", service.GetAllHandler)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", service.PostHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic("Listen and serve failed!")
	}
}
