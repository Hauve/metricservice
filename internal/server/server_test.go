package server

import (
	"github.com/Hauve/metricservice.git/internal/handlers"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	want := &MyServer{
		service: *handlers.New(),
		address: "localhost:8080",
	}
	if got := New(*handlers.New()); !reflect.DeepEqual(got, want) {
		t.Errorf("New() = %v, want %v", got, want)
	}
}

func Test_getAddress(t *testing.T) {
	tt := struct {
		name string
		want string
	}{
		name: "Test with default flag value",
		want: "localhost:8080",
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := getAddress(); got != tt.want {
			t.Errorf("getAddress() = %v, want %v", got, tt.want)
		}
	})

}
