package main

import (
	"net/http"

	"github.com/Hauve/metricservice.git/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()

	r.Get("/value/{metricType}/{metricName}", handlers.GetHandler)
	r.Get("/", handlers.GetAllHandler)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.PostHandler)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic("Listen and serve failed!")
	}
}
