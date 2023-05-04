package logger

import "net/http"

type (
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err == nil {
		r.WriteHeader(200)
	}
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	if statusCode != 200 {
		//Чтобы избежать предупреждение за повторную запись header.
		//Write по умолчанию выставляет его в 200, если запись ОК
		r.ResponseWriter.WriteHeader(statusCode)
	}
	r.responseData.status = statusCode
}
