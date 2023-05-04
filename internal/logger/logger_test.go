package logger

import (
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handlerForTest(w http.ResponseWriter, _ *http.Request) {
	data := []byte("body message")
	_, _ = w.Write(data)
}

func TestLogger_WithLogging(t *testing.T) {
	type fields struct {
		sugar zap.SugaredLogger
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "positive test",
			fields: fields{
				sugar: zap.SugaredLogger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := Logger{
				sugar: tt.fields.sugar,
			}
			_ = log.WithLogging(handlerForTest)
			t.SkipNow()
		})
	}
}

func Test_loggingResponseWriter_Write(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "ok",
			args: args{
				[]byte("for writing"),
			},
			wantStatus: 200,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &loggingResponseWriter{
				ResponseWriter: httptest.NewRecorder(),
				responseData:   &responseData{},
			}

			_, err := r.Write(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if r.responseData.status != tt.wantStatus {
				t.Errorf("Write() got status = %v, want status %v", r.responseData.status, tt.wantStatus)
			}
		})
	}
}

func Test_loggingResponseWriter_WriteHeader(t *testing.T) {
	type args struct {
		statusCode int
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "ok",
			args: args{
				200,
			},
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &loggingResponseWriter{
				ResponseWriter: httptest.NewRecorder(),
				responseData:   &responseData{},
			}
			r.WriteHeader(tt.args.statusCode)
			if r.responseData.status != tt.wantStatus {
				t.Errorf("Write() got status = %v, want status %v", r.responseData.status, tt.wantStatus)
			}
		})
	}
}
