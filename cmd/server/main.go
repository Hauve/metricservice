package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	data      map[string]string
	myMetrics []string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data:      make(map[string]string),
		myMetrics: []string{"gauge", "counter"},
	}
}

var MyMemStorage = NewMemStorage()

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/update/")
	pathParts := strings.Split(path, "/")

	ok := false
	for _, val := range MyMemStorage.myMetrics {
		if pathParts[0] == val {
			ok = true
			break
		}
	}
	if !ok {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	if len(pathParts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metric := pathParts[0]
	name := pathParts[1]
	val := pathParts[2]

	switch metric {
	case MyMemStorage.myMetrics[0]:
		_, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		MyMemStorage.data[name] = val
	case MyMemStorage.myMetrics[1]:
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
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update", mainHandler)
	mux.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotImplemented)
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
