package server

import (
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestService_GetHandler(t *testing.T) {
	type want struct {
		code  int
		value string
	}

	type fields struct {
		cfg     *config.ServerConfig
		storage storage.Storage
	}
	tests := []struct {
		name   string
		path   string
		fields fields
		want
	}{
		{
			name: "Positive test 1: get exists gauge",
			path: "/value/gauge/Key1",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "http://localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{200, "25.1"},
		},
		{
			name: "Positive test 2: get exists counter",
			path: "/value/counter/Key1",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "http://localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{200, "1"},
		},
		{
			name: "Negative test 1: get not exists gauge",
			path: "/value/gauge/Key2",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "http://localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{404, ""},
		},
		{
			name: "Negative test 2: get not exists counter",
			path: "/value/counter/Key2",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "http://localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{404, ""},
		},
		{
			name: "Negative test 3: get not exist metric type",
			path: "/value/socks/Key1",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "http://localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{404, ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MyServer{
				cfg:     tt.fields.cfg,
				storage: tt.fields.storage,
			}

			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, s.cfg.Address+tt.path, nil)

			s.storage.SetGauge("Key1", 25.1)
			s.storage.AddCounter("Key1", 1)
			s.GetHandler(resp, req)

			res := resp.Result()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				log.Printf("%e", err)
			}
			//static test gives warning with closing through defer
			err = res.Body.Close()
			if err != nil {
				log.Printf("%e", err)
			}

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.value, string(resBody))
			require.NoError(t, err)
			assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
		})
	}
}

func TestService_PostHandler(t *testing.T) {
	type want struct {
		code int
	}

	type fields struct {
		cfg     *config.ServerConfig
		storage storage.Storage
	}
	tests := []struct {
		name   string
		path   string
		fields fields
		want
	}{
		{
			name: "Positive test 1: post gauge",
			path: "/update/gauge/Key/55.6",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{http.StatusOK},
		},
		{
			name: "Positive test 2: post counter",
			path: "/update/counter/Key/94",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{http.StatusOK},
		},
		{
			name: "Negative test 1: more longer path",
			path: "/update/counter/Key/94/longer",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{http.StatusNotFound},
		},
		{
			name: "Negative test 2: shorter path",
			path: "/update/counter/Key",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{http.StatusNotFound},
		},
		{
			name: "Negative test 3: unknown metric type",
			path: "/update/counterCounter/Key/94",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{http.StatusNotImplemented},
		},
		{
			name: "Negative test 4: gauge, bad value",
			path: "/update/gauge/Key/sss",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{http.StatusBadRequest},
		},
		{
			name: "Negative test 5: counter, bad value",
			path: "/update/counter/Key/sss",
			fields: fields{
				cfg: &config.ServerConfig{
					Address: "localhost:8080",
				},
				storage: storage.NewMemStorage(),
			},
			want: want{http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			s := &MyServer{
				cfg:     tt.fields.cfg,
				storage: tt.fields.storage,
			}

			s.PostHandler(resp, req)
			res := resp.Result()

			assert.Equal(t, res.StatusCode, tt.want.code)
			assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")

			if tt.code == http.StatusOK {
				path := strings.Split(tt.path, "/")
				metricName := strings.ToLower(path[1])
				name := path[2]
				val := path[3]
				if metricName == "gauge" {
					resVal, ok := s.storage.GetGauge(name)
					assert.True(t, ok)

					flVal, _ := strconv.ParseFloat(val, 64)

					assert.Equal(t, resVal, flVal)
				} else if metricName == "counter" {
					resVal, ok := s.storage.GetCounter(name)
					assert.True(t, ok)

					intVal, _ := strconv.ParseInt(val, 10, 64)

					assert.Equal(t, resVal, intVal)
				}
			}

			//static test gives warning with closing through defer
			err := res.Body.Close()
			if err != nil {
				log.Printf("%e", err)
			}
		})
	}
}
