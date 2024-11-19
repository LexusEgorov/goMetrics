package runners

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	configPkg "github.com/LexusEgorov/goMetrics/internal/config"
	"github.com/LexusEgorov/goMetrics/internal/keeper"
	"github.com/LexusEgorov/goMetrics/internal/models"
	"github.com/LexusEgorov/goMetrics/internal/services/collectmetric"
	"github.com/LexusEgorov/goMetrics/internal/services/db"
	"github.com/LexusEgorov/goMetrics/internal/services/filestorage"
	storagePkg "github.com/LexusEgorov/goMetrics/internal/services/storage"
	"github.com/LexusEgorov/goMetrics/internal/transport"
)

type serverRunner struct{}

func (s serverRunner) Run(config configPkg.Server) error {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	sugar := logger.Sugar()
	router := chi.NewRouter()

	var storage keeper.Storager

	switch config.Mode {
	case configPkg.FileStorage:
		storage = filestorage.NewFileStorage(config.StorePath, config.StoreInterval, config.Restore)
	case configPkg.DBStorage:
		storage = db.NewDB(config.Host)
	default:
		storage = storagePkg.NewStorage(make(map[string]models.Metric))
	}

	keeper := keeper.NewKeeper(storage)
	transportServer := transport.NewServer(keeper, router, sugar)

	fmt.Println("Running server on", config.Host)
	fmt.Println("Backup interval: ", config.StoreInterval)
	fmt.Println("Backup file: ", config.StorePath)
	fmt.Println("Backup readed: ", config.Restore)
	fmt.Println("DB host: ", config.DB)
	fmt.Println("Storage mode: ", config.Mode)

	return http.ListenAndServe(config.Host, transportServer.Router)
}

func NewServer() *serverRunner {
	return &serverRunner{}
}

type agentRunner struct{}

func (a agentRunner) Run(config configPkg.Agent) {
	storage := storagePkg.NewStorage(make(map[string]models.Metric))
	keeper := keeper.NewKeeper(storage)

	var agent = collectmetric.NewAgent(config, keeper)

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
