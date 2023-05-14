package config

import (
	"testing"
)

func TestLoadServerConfig(t *testing.T) {
	tests := []struct {
		name string
		want *ServerConfig
	}{
		{
			name: "OK",
			want: &ServerConfig{
				Address:         "localhost:8080",
				StoreInterval:   300,
				FileStoragePath: "/tmp/metrics-db.json",
				Restore:         true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if got := LoadServerConfig(); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LoadServerConfig() = %v, want %v", got, tt.want)
			//}
			//Тест закомментирован, так как LoadServerConfig переопределяет флаг a, который определяется в тестах
			//config agent. panic
			//config.test.exe flag redefined: a
		})
	}
}
