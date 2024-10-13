package config

import (
	"flag"
	"os"
	"strconv"
)

type ServerVariables struct {
	Host string
}

func GetServerFlags() ServerVariables {
	var host string

	flag.StringVar(&host, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	envHost := os.Getenv("ADDRESS")

	if envHost != "" {
		host = envHost
	}

	return ServerVariables{
		Host: host,
	}
}

type AgentVariables struct {
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

func GetAgentFlags() AgentVariables {
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

	return AgentVariables{
		Host:           host,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}
}
