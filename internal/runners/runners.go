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
	"github.com/LexusEgorov/goMetrics/internal/models"
	"github.com/LexusEgorov/goMetrics/internal/services/collectmetric"
	"github.com/LexusEgorov/goMetrics/internal/services/filestorage"
	"github.com/LexusEgorov/goMetrics/internal/services/reader"
	"github.com/LexusEgorov/goMetrics/internal/services/saver"
	"github.com/LexusEgorov/goMetrics/internal/services/storage"
	"github.com/LexusEgorov/goMetrics/internal/transport"
)

type serverRunner struct{}

func (s serverRunner) Run(config config.Server) error {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	fileWriter := filestorage.NewFileWriter(config.StorePath)
	fileReader := filestorage.NewFileReader(config.StorePath)

	defer fileReader.Close()
	defer fileWriter.Close()

	initMetrics := make(map[string]models.Metric)

	if config.Restore {
		initMetrics = fileReader.Read()
	}

	saverRepo, readerRepo := storage.NewStorage(initMetrics)

	var metricSaver transport.Saver

	if config.StoreInterval == 0 {
		metricSaver = saver.NewSaver(saverRepo, fileWriter)
	} else {
		go fileWriter.RunSave(saverRepo, config.StoreInterval)
		metricSaver = saver.NewSaver(saverRepo, nil)
	}
	reader := reader.NewReader(readerRepo)

	sugar := logger.Sugar()
	router := chi.NewRouter()

	transportServer := transport.NewServer(metricSaver, reader, router, sugar)

	fmt.Println("Running server on", config.Host)
	fmt.Println("Backup interval: ", config.StoreInterval)
	fmt.Println("Backup file: ", config.StorePath)
	fmt.Println("Backup readed: ", config.Restore)
	return http.ListenAndServe(config.Host, transportServer.Router)
}

func NewServer() *serverRunner {
	return &serverRunner{}
}

type agentRunner struct{}

func (a agentRunner) Run(config config.Agent) {
	saverRepo, readerRepo := storage.NewStorage(make(map[string]models.Metric))

	saver := saver.NewSaver(saverRepo, nil)
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
