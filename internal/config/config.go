package config

import (
	"flag"
	"os"
	"strconv"
)

type server struct {
	Host          string
	StoreInterval int
	StorePath     string
	Restore       bool
}

func NewServer() server {
	var host string
	var storeInterval int
	var storePath string
	var restore bool

	flag.StringVar(&host, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&storePath, "i", "backup.txt", "store path")
	flag.IntVar(&storeInterval, "f", 300, "save interval")
	flag.BoolVar(&restore, "r", false, "is restore data?")
	flag.Parse()

	envHost := os.Getenv("ADDRESS")
	envInterval := os.Getenv("STORE_INTERVAL")
	envPath := os.Getenv("FILE_STORAGE_PATH")
	envRestore := os.Getenv("RESTORE")

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

	return server{
		Host:          host,
		StoreInterval: storeInterval,
		StorePath:     storePath,
		Restore:       restore,
	}
}

type Agent struct {
	Host           string
	ReportInterval int
	PollInterval   int
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
	var reportInterval int
	var pollInterval int

	flag.StringVar(&host, "a", "localhost:8080", "address and port for reporting")
	flag.IntVar(&reportInterval, "r", 10, "report interval")
	flag.IntVar(&pollInterval, "p", 2, "poll interval")
	flag.Parse()

	envHost := os.Getenv("ADDRESS")
	envReportInterval := os.Getenv("REPORT_INTERVAL")
	envPollInterval := os.Getenv("POLL_INTERVAL")

	if envHost != "" {
		host = envHost
	}

	if envReportInterval != "" {
		reportInterval = parseEnv(envReportInterval)
	}

	if envPollInterval != "" {
		pollInterval = parseEnv(envPollInterval)
	}

	return Agent{
		Host:           host,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}
}
