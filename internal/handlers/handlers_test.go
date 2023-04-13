package handlers

import (
	"github.com/Hauve/metricservice.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Service
	}{
		{
			name: "Positive test 1",
			want: &Service{
				MyMemStorage: storage.NewMemStorage(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetAllHandler(t *testing.T) {
	type fields struct {
		MyMemStorage Storage
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Positive test 1",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			s := &Service{
				MyMemStorage: tt.fields.MyMemStorage,
			}
			s.MyMemStorage.SetGauge("Key1", 25.1)
			s.MyMemStorage.SetCounter("Key1", 1)

			s.GetAllHandler(resp, req)

			res := resp.Result()
			defer func() {
				err := res.Body.Close()
				if err != nil {

				}
			}()

			testResultHTML := "<!DOCTYPE html><html><head><meta charset=\"UTF-8\"><title>Metrics</title></head><body>Gauge [keys] = values: <br>[Key1] = 25.1<br><br><br>Counters [keys] = values:<br>[Key1] = 1<br>\t</body>\n\t\t\t\t\t</html>"
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, string(resBody), testResultHTML)
			assert.Equal(t, res.Header.Get("Content-Type"), "text/html; charset=utf-8")
		})
	}
}

func TestService_GetHandler(t *testing.T) {
	type want struct {
		code  int
		value string
	}

	type fields struct {
		MyMemStorage Storage
	}
	tests := []struct {
		name   string
		path   string
		fields fields
		want
	}{
		{
			name: "Positive test 1: get exists gauge",
			path: "/value/gauge/Key1/",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{200, "25.1"},
		},
		{
			name: "Positive test 2: get exists counter",
			path: "/value/counter/Key1/",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{200, "1"},
		},
		{
			name: "Negative test 1: get not exists gauge",
			path: "/value/gauge/Key2/",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{404, ""},
		},
		{
			name: "Negative test 2: get not exists counter",
			path: "/value/counter/Key2/",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{404, ""},
		},
		{
			name: "Negative test 3: get not exist metric type",
			path: "/value/socks/Key1/",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{404, ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			s := &Service{
				MyMemStorage: tt.fields.MyMemStorage,
			}

			s.MyMemStorage.SetGauge("Key1", 25.1)
			s.MyMemStorage.SetCounter("Key1", 1)

			s.GetHandler(resp, req)

			res := resp.Result()
			defer func() {
				err := res.Body.Close()
				if err != nil {

				}
			}()
			resBody, err := io.ReadAll(res.Body)

			assert.Equal(t, res.StatusCode, tt.want.code)
			assert.Equal(t, string(resBody), tt.want.value)
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
		MyMemStorage Storage
	}
	tests := []struct {
		name   string
		path   string
		fields fields
		want
	}{
		{
			name: "Positive test 1: post gauge",
			path: "/update/gauge/Key/55.6/",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{200},
		},
		{
			name: "Positive test 2: post counter",
			path: "/update/counter/Key/94/",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{200},
		},
		{
			name: "Negative test 1: more longer path",
			path: "/update/counter/Key/94/longer",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{404},
		},
		{
			name: "Negative test 2: shorter path",
			path: "/update/counter/Key/",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{400},
		},
		{
			name: "Negative test 3: unknown metric type",
			path: "/update/counterCounter/Key/94/",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{404},
		},
		{
			name: "Negative test 4: gauge, bad value",
			path: "/update/gauge/Key/sss/",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{404},
		},
		{
			name: "Negative test 5: counter, bad value",
			path: "/update/counter/Key/sss/",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
			want: want{404},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			s := &Service{
				MyMemStorage: tt.fields.MyMemStorage,
			}

			s.PostHandler(resp, req)

			res := resp.Result()
			defer func() {
				err := res.Body.Close()
				if err != nil {

				}
			}()

			assert.Equal(t, res.StatusCode, tt.want.code)
			assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")

			path := strings.Split(tt.path, "/")
			metricName := strings.ToLower(path[1])
			name := path[2]
			val := path[3]

			if tt.code == http.StatusOK {
				if metricName == "gauge" {
					resVal, ok := s.MyMemStorage.GetGauge(name)
					assert.True(t, ok)

					flVal, _ := strconv.ParseFloat(val, 64)

					assert.Equal(t, resVal, flVal)
				} else if metricName == "counter" {
					resVal, ok := s.MyMemStorage.GetCounter(name)
					assert.True(t, ok)

					intVal, _ := strconv.ParseInt(val, 10, 64)

					assert.Equal(t, resVal, intVal)
				}
			}
		})
	}
}
