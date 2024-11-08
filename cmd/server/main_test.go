package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LexusEgorov/goMetrics/internal/services/reader"
	"github.com/LexusEgorov/goMetrics/internal/services/saver"
	"github.com/LexusEgorov/goMetrics/internal/services/storage"
	"github.com/LexusEgorov/goMetrics/internal/transport"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func setupRouter() *chi.Mux {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	saverRepo, readerRepo := storage.NewStorage()

	saver := saver.NewSaver(saverRepo)
	reader := reader.NewReader(readerRepo)

	sugar := logger.Sugar()
	router := chi.NewRouter()

	transportServer := transport.NewServer(saver, reader, router, sugar)

	return transportServer.Router
}

func Test_main(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		method   string
		url      string
		expected int
	}{
		{"GET", "/", http.StatusOK},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.url, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != test.expected {
			t.Errorf("Expected status %d, got %d for %s %s", test.expected, rr.Code, test.method, test.url)
		}
	}
}
