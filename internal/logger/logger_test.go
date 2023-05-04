package logger

import (
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"testing"
)

func TestLogger_WithLogging(t *testing.T) {
	type fields struct {
		sugar zap.SugaredLogger
	}
	type args struct {
		h http.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   http.HandlerFunc
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := Logger{
				sugar: tt.fields.sugar,
			}
			if got := log.WithLogging(tt.args.h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithLogging() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggingResponseWriter_Write(t *testing.T) {
	type fields struct {
		ResponseWriter http.ResponseWriter
		responseData   *responseData
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &loggingResponseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				responseData:   tt.fields.responseData,
			}
			got, err := r.Write(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggingResponseWriter_WriteHeader(t *testing.T) {
	type fields struct {
		ResponseWriter http.ResponseWriter
		responseData   *responseData
	}
	type args struct {
		statusCode int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &loggingResponseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				responseData:   tt.fields.responseData,
			}
			r.WriteHeader(tt.args.statusCode)
		})
	}
}
