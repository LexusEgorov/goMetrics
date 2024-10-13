package flags

import "flag"

type ServerFlags struct {
	Host string
}

func GetServerFlags() ServerFlags {
	var host string

	flag.StringVar(&host, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	return ServerFlags{
		Host: host,
	}
}

type AgentFlags struct {
	Host           string
	ReportInterval int
	PollInterval   int
}

func GetAgentFlags() AgentFlags {
	var host string
	var reportInterval int
	var pollInterval int

	flag.StringVar(&host, "a", "localhost:8080", "address and port for reporting")
	flag.IntVar(&reportInterval, "r", 10, "report interval")
	flag.IntVar(&pollInterval, "p", 2, "poll interval")
	flag.Parse()

	return AgentFlags{
		Host:           host,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}
}
