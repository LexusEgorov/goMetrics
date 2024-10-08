package storage

import (
	"reflect"
	"testing"
)

func TestCreateStorage(t *testing.T) {
	tests := []struct {
		name string
		want MemStorage
	}{
		{
			name: "Test create storage",
			want: MemStorage{
				data: make(map[MetricName]interface{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_AddGauge(t *testing.T) {
	type args struct {
		key   MetricName
		value Gauge
	}
	tests := []struct {
		name string
		m    *MemStorage
		args args
	}{
		{
			name: "Test AddGauge",
			m: &MemStorage{
				data: make(map[MetricName]interface{}),
			},
			args: args{
				key:   "Test",
				value: 1.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.AddGauge(tt.args.key, tt.args.value)

			if res := tt.m.Get(tt.args.key); res != tt.args.value {
				t.Errorf("AddGauge(%s) = %f, want %f", tt.args.key, res, tt.args.value)
			}
		})
	}
}
