package app

import (
	"flag"
	"fmt"
	"log"
	"marketflow/internal/api/handlers"
	"marketflow/internal/domain"
	"marketflow/internal/pkg/envzilla"
	"marketflow/internal/service"
	"net/http"
	"os"
	"strconv"
)

func init() {
	checkFlags()

	if err := envzilla.Loader("build/.env"); err != nil {
		log.Fatalf("Config file load error: %s", err.Error())
	}
}

// Setup function sets connection to the adapters
func Setup(db domain.Database, cacheMemory domain.CacheMemory, datafetcher domain.DataFetcher) *http.ServeMux {
	datafetchServ := service.NewDataFetcher(datafetcher, db, cacheMemory)

	modeHandler := handlers.NewSwitchModeHandler(datafetchServ)
	marketHandler := handlers.NewMarketDataHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /mode/{mode}", modeHandler.SwitchMode)

	mux.HandleFunc("GET /health", modeHandler.CheckHealth)

	mux.HandleFunc("GET /prices/{metric}/{symbol}", marketHandler.ProcessMetricQueryByAll)
	mux.HandleFunc("GET /prices/{metric}/{exchange}/{symbol}", marketHandler.ProcessMetricQueryByExchange)

	return mux
}

// checkFlags validate CLI flags
func checkFlags() {
	flag.Parse()
	portNum, err := strconv.Atoi(*domain.Port)
	if err != nil {
		log.Fatalf("Port number is incorrect: %s ", *domain.Port)
	}

	if portNum < 1024 || portNum > 65535 {
		log.Fatalf("Port number is incorrect: %d , must be in range 1024 and 65535 ", portNum)
	}

	if *domain.HelpFlag {
		printHelp()
	}

}

// Prints help message
func printHelp() {
	fmt.Println(`Usage:
  marketflow [--port <N>]
  marketflow --help

Options:
  --port N     Port number`)
	os.Exit(0)
}
