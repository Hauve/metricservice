package server

import (
	"bytes"
	"encoding/json"
	"github.com/Hauve/metricservice.git/internal/config"
	"github.com/Hauve/metricservice.git/internal/jsonmodel"
	"github.com/Hauve/metricservice.git/internal/logger"
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/go-chi/chi/v5"
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

func TestMyServer_GetHandler(t *testing.T) {
	type want struct {
		code  int
		value string
	}
	tests := []struct {
		name string
		path string
		want
	}{
		{
			name: "Positive test 1: get exists gauge",
			path: "/value/gauge/Key1",
			want: want{200, "25.1"},
		},
		{
			name: "Positive test 2: get exists counter",
			path: "/value/counter/Key1",
			want: want{200, "1"},
		},
		{
			name: "Negative test 1: get not exists gauge",
			path: "/value/gauge/Key2",
			want: want{404, ""},
		},
		{
			name: "Negative test 2: get not exists counter",
			path: "/value/counter/Key2",
			want: want{404, ""},
		},
		{
			name: "Negative test 3: get not exist metric type",
			path: "/value/socks/Key1",
			want: want{404, ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lg, err := logger.New()
			require.NoError(t, err)

			s := &MyServer{
				cfg: &config.ServerConfig{
					Address: "http://localhost:8080",
				},
				storage: storage.NewMemStorage(),
				router:  chi.NewRouter(),
				logger:  *lg,
			}
			s.storage.SetGauge("Key1", 25.1)
			s.storage.AddCounter("Key1", 1)

			s.registerRoutes()

			req := httptest.NewRequest(http.MethodGet, s.cfg.Address+tt.path, nil)
			resp := httptest.NewRecorder()
			s.router.ServeHTTP(resp, req)

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

func TestMyServer_PostHandler(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		path string
		want
	}{
		{
			name: "Positive test 1: post gauge",
			path: "/update/gauge/Key/55.6",
			want: want{http.StatusOK},
		},
		{
			name: "Positive test 2: post counter",
			path: "/update/counter/Key/94",
			want: want{http.StatusOK},
		},
		{
			name: "Negative test 1: more longer path",
			path: "/update/counter/Key/94/longer",
			want: want{http.StatusNotFound},
		},
		{
			name: "Negative test 2: shorter path",
			path: "/update/counter/Key",
			want: want{http.StatusNotFound},
		},
		{
			name: "Negative test 3: unknown metric type",
			path: "/update/counterCounter/Key/94",
			want: want{http.StatusNotImplemented},
		},
		{
			name: "Negative test 4: gauge, bad value",
			path: "/update/gauge/Key/sss",
			want: want{http.StatusBadRequest},
		},
		{
			name: "Negative test 5: counter, bad value",
			path: "/update/counter/Key/sss",
			want: want{http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lg, err := logger.New()
			require.NoError(t, err)

			s := &MyServer{
				cfg: &config.ServerConfig{
					Address: "http://localhost:8080",
				},
				storage: storage.NewMemStorage(),
				router:  chi.NewRouter(),
				logger:  *lg,
			}

			s.registerRoutes()

			req := httptest.NewRequest(http.MethodPost, s.cfg.Address+tt.path, nil)
			resp := httptest.NewRecorder()
			s.router.ServeHTTP(resp, req)

			res := resp.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(
				t,
				"text/plain; charset=utf-8",
				res.Header.Get("Content-Type"),
			)

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
			err = res.Body.Close()
			if err != nil {
				log.Printf("%e", err)
			}
		})
	}
}

func TestMyServer_JSONGetHandler(t *testing.T) {

	tests := []struct {
		name string
		code int
		body jsonmodel.Metrics
	}{
		{
			name: "positive test 1: get gauge",
			code: http.StatusOK,
			body: jsonmodel.Metrics{
				ID:    "name1",
				MType: "gauge",
			},
		},
		{
			name: "positive test 2: get counter",
			code: http.StatusOK,
			body: jsonmodel.Metrics{
				ID:    "name1",
				MType: "counter",
			},
		},
		//{
		//	name: "negative test 1: get non-existent metric",
		//	code: http.StatusNotFound,
		//	body: jsonmodel.Metrics{
		//		ID:    "name",
		//		MType: "counter",
		//	},
		//},
		//{
		//	name: "negative test 2: get non-existent metric type",
		//	code: http.StatusNotFound,
		//	body: jsonmodel.Metrics{
		//		ID:    "name1",
		//		MType: "non-existent type",
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MyServer{
				storage: storage.NewMemStorage(),
				router:  chi.NewRouter(),
			}
			s.storage.SetGauge("name1", 25.1)
			s.storage.SetCounter("name1", 25)

			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Errorf("json marshalling failed: %s", err)
			}

			buf := bytes.NewBuffer(body)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/value/", buf)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			s.router.Post("/value/", s.JSONGetHandler)
			s.router.ServeHTTP(resp, req)

			res := resp.Result()

			require.Equal(t, tt.code, res.StatusCode)
			assert.Equal(
				t,
				"application/json; charset=utf-8",
				res.Header.Get("Content-Type"),
			)
			body, err = io.ReadAll(res.Body)
			assert.NoError(t, err)

			respData := jsonmodel.Metrics{}
			err = json.Unmarshal(body, &respData)
			assert.NoError(t, err)

			assert.Equal(t, tt.body.ID, respData.ID)
			assert.Equal(t, tt.body.MType, respData.MType)
			if tt.body.MType == storage.Gauge {
				assert.Equal(t, 25.1, *respData.Value)
			} else if tt.body.MType == storage.Counter {
				assert.Equal(t, int64(25), *respData.Delta)
			}
		})
	}
}

func TestMyServer_JSONPostHandler(t *testing.T) {
	type fields struct {
		cfg     *config.ServerConfig
		storage storage.Storage
		router  chi.Router
		logger  logger.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MyServer{
				cfg:     tt.fields.cfg,
				storage: tt.fields.storage,
				router:  tt.fields.router,
				logger:  tt.fields.logger,
			}
			s.JSONPostHandler(tt.args.w, tt.args.r)
		})
	}
}
