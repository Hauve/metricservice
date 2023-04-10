package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name   string
		path   string
		method string
		want   want
	}{
		{
			name:   "positive test #1",
			path:   "/update/gauge/MetricName/255.6",
			method: http.MethodPost,
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		}, {
			name:   "positive test #2",
			path:   "/update/counter/counterName/54",
			method: http.MethodPost,
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		}, {
			name:   "negative test #1",
			path:   "/update/counterCounterCounter/counterName/54",
			method: http.MethodPost,
			want: want{
				code:        501,
				contentType: "text/plain; charset=utf-8",
			},
		}, {
			name:   "negative test #2",
			path:   "/update/counter/counterName/elf",
			method: http.MethodPost,
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		}, {
			name:   "negative test #3",
			path:   "/update/counter/counterName/potato",
			method: http.MethodPost,
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		}, {
			name:   "negative test #4",
			path:   "/counter/counterName",
			method: http.MethodPost,
			want: want{
				code:        501,
				contentType: "text/plain; charset=utf-8",
			},
		}, {
			name:   "negative test #5",
			path:   "/update/counter/counterName/54",
			method: http.MethodPatch,
			want: want{
				code:        405,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			PostHandler(w, request)

			res := w.Result()
			err := res.Body.Close()
			if err != nil {
				return
			}

			assert.Equal(t, res.StatusCode, test.want.code)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}
