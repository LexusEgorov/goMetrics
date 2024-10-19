package transport

import (
	"net/http"
	"reflect"
	"testing"
)

func TestTransportLayer_UpdateMetric(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		tr   transportLayer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.UpdateMetric(tt.args.w, tt.args.r)
		})
	}
}

func TestTransportLayer_GetMetric(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		tr   transportLayer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.GetMetric(tt.args.w, tt.args.r)
		})
	}
}

func TestTransportLayer_GetMetrics(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		tr   transportLayer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.GetMetrics(tt.args.w, tt.args.r)
		})
	}
}

func TestTransportLayer_SendMetric(t *testing.T) {
	type args struct {
		metricName  string
		metricType  string
		metricValue string
	}
	tests := []struct {
		name string
		tr   transportLayer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.SendMetric("localhost:8080", tt.args.metricName, tt.args.metricType, tt.args.metricValue)
		})
	}
}

func TestCreateTransport(t *testing.T) {
	tests := []struct {
		name string
		want transportLayer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateTransport(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTransport() = %v, want %v", got, tt.want)
			}
		})
	}
}
