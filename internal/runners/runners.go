package runners

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/LexusEgorov/goMetrics/internal/config"
	"github.com/LexusEgorov/goMetrics/internal/services/collectmetric"
	"github.com/LexusEgorov/goMetrics/internal/services/reader"
	"github.com/LexusEgorov/goMetrics/internal/services/saver"
	"github.com/LexusEgorov/goMetrics/internal/services/storage"
	"github.com/LexusEgorov/goMetrics/internal/transport"
)

type serverRunner struct{}

func (s serverRunner) Run(host string) error {
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

	fmt.Println("Running server on", host)
	return http.ListenAndServe(host, transportServer.Router)
}

func NewServer() *serverRunner {
	return &serverRunner{}
}

type agentRunner struct{}

func (a agentRunner) Run(config config.Agent) {
	saverRepo, readerRepo := storage.NewStorage()

	saver := saver.NewSaver(saverRepo)
	reader := reader.NewReader(readerRepo)
	var agent = collectmetric.NewAgent(config, saver, reader)

	stopChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		close(stopChan)
	}()

	agent.Start(stopChan)
}

func NewAgent() *agentRunner {
	return &agentRunner{}
}
