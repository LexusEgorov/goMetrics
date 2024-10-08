package collectmetric

import (
	"reflect"
	"testing"

	"github.com/LexusEgorov/goMetrics/internal/services/storage"
)

func TestCreateAgent(t *testing.T) {
	tests := []struct {
		name string
		want MetricsCollector
	}{
		{
			name: "DefaultAgent",
			want: &MetricAgent{
				storage:   storage.CreateStorage(),
				pollCount: 0,
				intervals: agentIntervals{
					collect: 2,
					send:    10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateAgent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}
