package logger

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Logger struct {
	zap.SugaredLogger
}

func New() (*Logger, error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("logger creating failed: %w", err)
	}
	sug := log.Sugar()
	return &Logger{
		SugaredLogger: *sug,
	}, nil
}

func (log Logger) WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   rData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", rData.status,
			"duration", duration*time.Millisecond,
			"size", rData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
