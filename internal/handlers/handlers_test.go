package handlers

import (
	"github.com/Hauve/metricservice.git/internal/storage"
	"net/http"
	"reflect"
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
	type args struct {
		w   http.ResponseWriter
		in1 *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
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
			s := &Service{
				MyMemStorage: tt.fields.MyMemStorage,
			}
			s.GetAllHandler(tt.args.w, tt.args.in1)
		})
	}
}

func TestService_GetHandler(t *testing.T) {
	type fields struct {
		MyMemStorage Storage
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
		{
			name: "Positive test 1",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				MyMemStorage: tt.fields.MyMemStorage,
			}
			s.GetHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestService_PostHandler(t *testing.T) {
	type fields struct {
		MyMemStorage Storage
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
		{
			name: "Positive test 1",
			fields: fields{
				MyMemStorage: storage.NewMemStorage(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				MyMemStorage: tt.fields.MyMemStorage,
			}
			s.PostHandler(tt.args.w, tt.args.r)
		})
	}
}
