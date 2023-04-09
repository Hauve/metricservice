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

//Error:      	Not equal:
//expected: 404
//actual  : 501
//Messages:   	Несоответствие статус кода ответа ожидаемому в хендлере 'POST http://localhost:8080/updater/counter/testCounter/100'
//iteration1_test.go:242: Оригинальный запрос:
//POST /updater/counter/testCounter/100 HTTP/1.1

//Error:      	Not equal:
//expected: 501
//actual  : 404
//Messages:   	Несоответствие статус кода ответа ожидаемому в хендлере 'POST http://localhost:8080/update/unknown/testCounter/100'
//iteration1_test.go:227: Оригинальный запрос:
//POST /update/unknown/testCounter/100 HTTP/1.1
