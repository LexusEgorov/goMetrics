// Пакет config определяет переменные для старта приложения.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Режимы хранилища.
const (
	MemStorage = iota
	FileStorage
	DBStorage
)

// Переменные для запуска сервера.
type Server struct {
	Host          string `json:"address"`
	StoreInterval int    `json:"store_interval"`
	StorePath     string `json:"store_file"`
	Restore       bool   `json:"restore"`
	DB            string `json:"database_dsn"`
	Key           string `json:"crypto_key"`
	Mode          int
}

type ServerJSON struct {
	Server
	StoreInterval string `json:"store_interval"`
}

func parseServerJSON(configBytes []byte) *Server {
	if len(configBytes) == 0 {
		return &Server{}
	}

	serverJSON := ServerJSON{}
	err := json.Unmarshal(configBytes, &serverJSON)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &Server{
		Host:          serverJSON.Host,
		StoreInterval: parseTime(serverJSON.StoreInterval),
		StorePath:     serverJSON.StorePath,
		Restore:       serverJSON.Restore,
		DB:            serverJSON.DB,
		Key:           serverJSON.Key,
	}
}

// NewServer определяет переменные из флагов командной строки и переменных окружения для сервера.
func NewServer() Server {
	mode := MemStorage

	var host string
	var storeInterval int
	var storePath string
	var restore bool
	var db string
	var key string
	var cryptoKey string

	var configJSON string

	//Парсинг флагов командной строки.
	flag.StringVar(&host, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&storePath, "f", "backup.txt", "store path")
	flag.StringVar(&db, "d", "", "db path")
	flag.StringVar(&key, "k", "", "secret key")
	flag.StringVar(&cryptoKey, "crypto-key", "", "crypto key")
	flag.StringVar(&configJSON, "c", "", "JSON config")
	flag.IntVar(&storeInterval, "i", 300, "save interval")
	flag.BoolVar(&restore, "r", false, "is restore data?")
	flag.Parse()

	JSONbytes, err := os.ReadFile(configJSON)

	if err != nil {
		fmt.Println(err)
		JSONbytes = []byte{}
	}

	ServerJSON := parseServerJSON(JSONbytes)

	if host == "" {
		host = ServerJSON.Host
	}

	if storeInterval == 0 {
		storeInterval = ServerJSON.StoreInterval
	}

	if storePath == "" {
		storePath = ServerJSON.StorePath
	}

	if check := flag.Lookup("r"); check == nil {
		restore = ServerJSON.Restore
	}

	if db == "" {
		db = ServerJSON.DB
	}

	if cryptoKey == "" {
		cryptoKey = ServerJSON.Key
	}

	//Получение переменных окружения.
	envHost := os.Getenv("ADDRESS")
	envInterval := os.Getenv("STORE_INTERVAL")
	envPath := os.Getenv("FileStorage_PATH")
	envRestore := os.Getenv("RESTORE")
	envDB := os.Getenv("DATABASE_DSN")
	envKey := os.Getenv("KEY")
	envCryptoKey := os.Getenv("CRYPTO_KEY")

	if envKey != "" {
		key = envKey
	}

	if envCryptoKey != "" {
		cryptoKey = envCryptoKey
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

	if cryptoKey != "" {
		key = cryptoKey
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

// Переменные для запуска сбора метрик.
type Agent struct {
	Host           string `json:"address"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
	Key            string `json:"crypto_key"`
	RateLimit      int
}

type AgentJSON struct {
	Agent
	ReportInterval string `json:"report_interval"`
	PollInterval   string `json:"poll_interval"`
}

func parseEnv(variable string) int {
	parsed, err := strconv.ParseInt(variable, 0, 64)

	if err != nil {
		panic(err)
	}

	return int(parsed)
}

func parseTime(time string) int {
	numPart := time[:len(time)-1]
	unitPart := time[len(time)-1:]

	num, err := strconv.Atoi(numPart)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	switch unitPart {
	case "s":
		return num
	case "m":
		return num * 60
	case "h":
		return num * 3600
	default:
		fmt.Println("Unsupported time (JSON config)")
		return 0
	}
}

func parseAgentJSON(configBytes []byte) *Agent {
	if len(configBytes) == 0 {
		return &Agent{}
	}

	agentJSON := AgentJSON{}
	err := json.Unmarshal(configBytes, &agentJSON)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &Agent{
		Host:           agentJSON.Host,
		ReportInterval: parseTime(agentJSON.ReportInterval),
		PollInterval:   parseTime(agentJSON.PollInterval),
		Key:            agentJSON.Key,
	}
}

// NewAgent определяет переменные из флагов командной строки и переменных окружения для агента.
func NewAgent() Agent {
	var host string
	var key string
	var cryptoKey string
	var reportInterval int
	var pollInterval int
	var rateLimit int

	var configJSON string

	//Парсинг флагов командной строки.
	flag.StringVar(&host, "a", "localhost:8080", "address and port for reporting")
	flag.StringVar(&key, "k", "", "secret key")
	flag.StringVar(&cryptoKey, "crypto-key", "", "crypto key")
	flag.StringVar(&configJSON, "c", "", "JSON config")
	flag.IntVar(&reportInterval, "r", 10, "report interval")
	flag.IntVar(&pollInterval, "p", 2, "poll interval")
	flag.IntVar(&rateLimit, "l", 1, "rate limit")
	flag.Parse()

	JSONbytes, err := os.ReadFile(configJSON)

	if err != nil {
		fmt.Println(err)
		JSONbytes = []byte{}
	}

	AgentJSON := parseAgentJSON(JSONbytes)

	if host == "" {
		host = AgentJSON.Host
	}

	if cryptoKey == "" {
		cryptoKey = AgentJSON.Key
	}

	if pollInterval == 0 {
		pollInterval = AgentJSON.PollInterval
	}

	if reportInterval == 0 {
		reportInterval = AgentJSON.ReportInterval
	}

	//Получение переменных окружения.
	envHost := os.Getenv("ADDRESS")
	envReportInterval := os.Getenv("REPORT_INTERVAL")
	envPollInterval := os.Getenv("POLL_INTERVAL")
	envKey := os.Getenv("KEY")
	envCryptoKey := os.Getenv("CRYPTO_KEY")
	envLimit := os.Getenv("RATE_LIMIT")

	if envKey != "" {
		key = envKey
	}

	if envCryptoKey != "" {
		cryptoKey = envCryptoKey
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

	if cryptoKey != "" {
		key = cryptoKey
	}

	return Agent{
		Host:           host,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
		Key:            key,
		RateLimit:      rateLimit,
	}
}
