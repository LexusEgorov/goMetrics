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
	"github.com/LexusEgorov/goMetrics/internal/services/sign"
	storagePkg "github.com/LexusEgorov/goMetrics/internal/services/storage"
	"github.com/LexusEgorov/goMetrics/internal/transport"
)

type serverRunner struct{}

func (s serverRunner) Run(config configPkg.Server, stopChan chan struct{}) error {
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
		storage = db.NewDB(config.DB)
	default:
		storage = storagePkg.NewStorage(make(map[string]models.Metric))
	}

	defer storage.Close()

	keeper := keeper.NewKeeper(storage)
	signer := sign.NewSign(config.Key)

	transportServer := transport.NewServer(keeper, router, sugar, signer)

	fmt.Println("Running server on", config.Host)
	fmt.Println("Backup interval: ", config.StoreInterval)
	fmt.Println("Backup file: ", config.StorePath)
	fmt.Println("Backup readed: ", config.Restore)
	fmt.Println("DB host: ", config.DB)
	fmt.Println("KEY: ", config.Key)
	fmt.Println("Storage mode: ", config.Mode)

	go func() {
		shutdown := false

		for !shutdown {
			select {
			case <-stopChan:
				shutdown = true
				fmt.Println("shutting down")
				transportServer.Shutdown()
				fmt.Println("shutted down")
				return
			default:
			}
		}
	}()

	return http.ListenAndServe(config.Host, transportServer.Router)
}

func NewServer() *serverRunner {
	return &serverRunner{}
}

type agentRunner struct{}

func (a agentRunner) Run(config configPkg.Agent) {
	signer := sign.NewSign(config.Key)

	var agent = collectmetric.NewAgent(config, signer)

	stopChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-signalChan
		close(stopChan)
	}()

	agent.Start(stopChan)
}

func NewAgent() *agentRunner {
	return &agentRunner{}
}
