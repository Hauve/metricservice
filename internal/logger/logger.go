package logger

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Logger struct {
	sugar zap.SugaredLogger
}

func New() (*Logger, error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		err = fmt.Errorf("logger creating failed: %w", err)
		return nil, err
	}
	sug := log.Sugar()
	return &Logger{
		sugar: *sug,
	}, nil
}

func (log Logger) WithLogging(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		fmt.Println("Before ServeHTTP in Logger")
		h.ServeHTTP(&lw, r)
		fmt.Println("After ServeHTTP in Logger")

		duration := time.Since(start)

		log.sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration*time.Millisecond,
			"size", responseData.size,
		)
	}
	return logFn
}
