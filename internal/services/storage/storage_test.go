package storage

// import (
// 	"reflect"
// 	"testing"
// )

// func TestNewStorage(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want memStorage
// 	}{
// 		{
// 			name: "Test create storage",
// 			want: memStorage{
// 				data: make(map[string]Metric),
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			storage := NewStorage()

// 			got := storage.GetAll()
// 			t.Logf("Got: %v, Want: %v", got, tt.want.data)

// 			if !reflect.DeepEqual(got, tt.want.data) {
// 				t.Errorf("NewStorage().GetAll() = %v, want %v", got, tt.want.data)
// 			}
// 		})
// 	}
// }

// func TestMemStorage_AddGauge(t *testing.T) {
// 	type args struct {
// 		key   string
// 		value float64
// 	}
// 	tests := []struct {
// 		name string
// 		m    *memStorage
// 		args args
// 	}{
// 		{
// 			name: "Test AddGauge",
// 			m: &memStorage{
// 				data: make(map[string]Metric),
// 			},
// 			args: args{
// 				key:   "Test",
// 				value: 1.0,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.m.AddGauge(tt.args.key, tt.args.value)

// 			// if res := tt.m.Get(tt.args.key); res != tt.args.value {
// 			// 	t.Errorf("AddGauge(%s) = %f, want %f", tt.args.key, res, tt.args.value)
// 			// }
// 		})
// 	}
// }

// func TestMemStorage_GetAll(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		m    memStorage
// 		want map[string]Metric
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.m.GetAll(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("memStorage.GetAll() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestMemStorage_GetCounter(t *testing.T) {
// 	type args struct {
// 		key string
// 	}
// 	tests := []struct {
// 		name  string
// 		m     memStorage
// 		args  args
// 		want  int64
// 		want1 bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, got1 := tt.m.GetCounter(tt.args.key)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("memStorage.GetCounter() got = %v, want %v", got, tt.want)
// 			}
// 			if got1 != tt.want1 {
// 				t.Errorf("memStorage.GetCounter() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }

// func TestMemStorage_GetGauge(t *testing.T) {
// 	type args struct {
// 		key string
// 	}
// 	tests := []struct {
// 		name  string
// 		m     memStorage
// 		args  args
// 		want  float64
// 		want1 bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, got1 := tt.m.GetGauge(tt.args.key)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("memStorage.GetGauge() got = %v, want %v", got, tt.want)
// 			}
// 			if got1 != tt.want1 {
// 				t.Errorf("memStorage.GetGauge() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }

// func TestMemStorage_AddCounter(t *testing.T) {
// 	type args struct {
// 		key   string
// 		value int64
// 	}
// 	tests := []struct {
// 		name string
// 		m    *memStorage
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.m.AddCounter(tt.args.key, tt.args.value)
// 		})
// 	}
// }
