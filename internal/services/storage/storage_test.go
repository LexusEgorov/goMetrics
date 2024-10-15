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
			storage := CreateStorage()

			got := storage.GetAll()
			t.Logf("Got: %v, Want: %v", got, tt.want.data)

			if !reflect.DeepEqual(got, tt.want.data) {
				t.Errorf("CreateStorage().GetAll() = %v, want %v", got, tt.want.data)
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

			// if res := tt.m.Get(tt.args.key); res != tt.args.value {
			// 	t.Errorf("AddGauge(%s) = %f, want %f", tt.args.key, res, tt.args.value)
			// }
		})
	}
}

func TestMemStorage_GetAll(t *testing.T) {
	tests := []struct {
		name string
		m    MemStorage
		want map[MetricName]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type args struct {
		key MetricName
	}
	tests := []struct {
		name  string
		m     MemStorage
		args  args
		want  Counter
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.m.GetCounter(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetCounter() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MemStorage.GetCounter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type args struct {
		key MetricName
	}
	tests := []struct {
		name  string
		m     MemStorage
		args  args
		want  Gauge
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.m.GetGauge(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetGauge() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MemStorage.GetGauge() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	type args struct {
		key   MetricName
		value Counter
	}
	tests := []struct {
		name string
		m    *MemStorage
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.AddCounter(tt.args.key, tt.args.value)
		})
	}
}

func TestCounter_String(t *testing.T) {
	tests := []struct {
		name string
		c    Counter
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Counter.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGauge_String(t *testing.T) {
	tests := []struct {
		name string
		g    Gauge
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.String(); got != tt.want {
				t.Errorf("Gauge.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
