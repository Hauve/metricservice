package main

import (
	"net/http"

	"github.com/Hauve/metricservice.git/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.UpdateHandler)
	mux.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
	}))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic("Listen and serve failed!")
	}
}
