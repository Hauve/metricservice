package logger

import "net/http"

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	if err == nil {
		r.WriteHeader(200)
	}
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	if statusCode != 200 {
		//Чтобы избежать предупреждение за повторную запись header.
		//Write по умолчанию выставляет его в 200, если запись ОК
		r.ResponseWriter.WriteHeader(statusCode)
	}
	r.responseData.status = statusCode // захватываем код статуса
}
