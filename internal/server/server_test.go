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
