package config

import (
	"flag"
	"os"
	"strconv"
)

const (
	MemStorage = iota
	FileStorage
	DBStorage
)

type Server struct {
	Host          string
	StoreInterval int
	StorePath     string
	Restore       bool
	DB            string
	Key           string
	Mode          int
}

func NewServer() Server {
	mode := MemStorage

	var host string
	var storeInterval int
	var storePath string
	var restore bool
	var db string
	var key string

	flag.StringVar(&host, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&storePath, "i", "backup.txt", "store path")
	flag.StringVar(&db, "d", "", "db path")
	flag.StringVar(&key, "k", "", "secret key")
	flag.IntVar(&storeInterval, "f", 300, "save interval")
	flag.BoolVar(&restore, "r", false, "is restore data?")
	flag.Parse()

	envHost := os.Getenv("ADDRESS")
	envInterval := os.Getenv("STORE_INTERVAL")
	envPath := os.Getenv("FileStorage_PATH")
	envRestore := os.Getenv("RESTORE")
	envDB := os.Getenv("DATABASE_DSN")
	envKey := os.Getenv("KEY")

	if envKey != "" {
		key = envKey
	}

	if envHost != "" {
		host = envHost
	}

	if envInterval != "" {
		parsedInterval, err := strconv.Atoi(envInterval)

		if err == nil {
			storeInterval = parsedInterval
		}
	}

	if envPath != "" {
		storePath = envPath
	}

	if envRestore != "" {
		parsedRestore, err := strconv.ParseBool(envRestore)

		if err == nil {
			restore = parsedRestore
		}
	}

	if envDB != "" {
		db = envDB
	}

	if storePath != "" {
		mode = FileStorage
	}

	if db != "" {
		mode = DBStorage
	}

	return Server{
		Host:          host,
		StoreInterval: storeInterval,
		StorePath:     storePath,
		Restore:       restore,
		DB:            db,
		Key:           key,
		Mode:          mode,
	}
}

type Agent struct {
	Host           string
	ReportInterval int
	PollInterval   int
	Key            string
	RateLimit      int
}

func parseEnv(variable string) int {
	parsed, err := strconv.ParseInt(variable, 0, 64)

	if err != nil {
		panic(err)
	}

	return int(parsed)
}

func NewAgent() Agent {
	var host string
	var key string
	var reportInterval int
	var pollInterval int
	var rateLimit int

	flag.StringVar(&host, "a", "localhost:8080", "address and port for reporting")
	flag.StringVar(&key, "k", "", "secret key")
	flag.IntVar(&reportInterval, "r", 10, "report interval")
	flag.IntVar(&pollInterval, "p", 2, "poll interval")
	flag.IntVar(&rateLimit, "l", 1, "rate limit")
	flag.Parse()

	envHost := os.Getenv("ADDRESS")
	envReportInterval := os.Getenv("REPORT_INTERVAL")
	envPollInterval := os.Getenv("POLL_INTERVAL")
	envKey := os.Getenv("KEY")
	envLimit := os.Getenv("RATE_LIMIT")

	if envKey != "" {
		key = envKey
	}

	if envHost != "" {
		host = envHost
	}

	if envReportInterval != "" {
		reportInterval = parseEnv(envReportInterval)
	}

	if envPollInterval != "" {
		pollInterval = parseEnv(envPollInterval)
	}

	if envLimit != "" {
		rateLimit = parseEnv(envLimit)
	}

	return Agent{
		Host:           host,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
		Key:            key,
		RateLimit:      rateLimit,
	}
}
