package sender

import (
	"encoding/binary"
	"testing"
)

func Test_compress(t *testing.T) {
	tests := []struct {
		name    string
		arg     []byte
		wantErr bool
	}{
		{
			name: "OK",
			arg: []byte("Next, we define a mock controller " +
				"inside our test. A mock controller is responsible " +
				"for tracking and asserting the expectations of its " +
				"associated mock objects."),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compress(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if binary.Size(got) > binary.Size(tt.arg) {
				t.Errorf("compress() gotten length %d > argument length %d", len(got), len(tt.arg))
			}
		})
	}
}
