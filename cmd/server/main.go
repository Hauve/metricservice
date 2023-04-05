package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	data map[string]string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]string),
	}
}

var MyMemStorage = NewMemStorage()

func gaugeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/update/gauge/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	name := pathParts[0]
	val := pathParts[1]
	_, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	MyMemStorage.data[name] = val
	w.WriteHeader(http.StatusOK)

	//_, err = w.Write([]byte("OK"))
	//if err != nil {
	//	panic("Недоступна запись в response writer")
	//}
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/update/counter/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	name := pathParts[0]
	val := pathParts[1]
	floatNewValue, err := strconv.ParseFloat(val, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, ok := MyMemStorage.data[name]; !ok {
		MyMemStorage.data[name] = val
	} else {
		oldValue := MyMemStorage.data[name]
		floatOldValue, _ := strconv.ParseFloat(oldValue, 64)
		newValue := fmt.Sprintf("%f", floatNewValue+floatOldValue)
		MyMemStorage.data[name] = newValue
	}
	w.WriteHeader(http.StatusOK)

	//_, err = w.Write([]byte("OK"))
	//if err != nil {
	//	panic("Недоступна запись в response writer")
	//}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/gauge/", gaugeHandler)
	mux.HandleFunc("/update/counter/", counterHandler)
	mux.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotFound)
	}))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic("Listen and serve failed!")
	}
}
