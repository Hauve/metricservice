package server

import (
	"flag"
	"github.com/Hauve/metricservice.git/internal/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

type MyServer struct {
	service handlers.Service

	address string
}

func New(service handlers.Service) *MyServer {
	var address *string
	addrEnv, ok := os.LookupEnv("ADDRESS")
	if !ok {
		address = flag.String("a", "localhost:8080", "address")
		flag.Parse()
	} else {
		address = &addrEnv
	}

	log.Println(*address)
	return &MyServer{
		service: service,
		address: *address,
	}
}

func (serv *MyServer) Run() {
	service := handlers.New()

	r := chi.NewRouter()

	r.Get("/value/{metricType}/{metricName}", service.GetHandler)
	r.Get("/", service.GetAllHandler)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", service.PostHandler)

	err := http.ListenAndServe(serv.address, r)
	if err != nil {
		panic("Listen and serve failed!")
	}
}
