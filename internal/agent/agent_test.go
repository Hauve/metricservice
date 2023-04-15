package agent

import (
	"github.com/Hauve/metricservice.git/internal/storage"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tt := struct {
		name string
		want *MyAgent
	}{
		name: "Single test New",
		want: &MyAgent{
			storage:        storage.NewMemStorage(),
			address:        "localhost:8080",
			reportInterval: 10,
			pollInterval:   2,
		},
	}
	t.Run(tt.name, func(t *testing.T) {
		if got := New(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("New() = %v, want %v", got, tt.want)
		}
	})
}
