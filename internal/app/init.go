package app

import (
	"flag"
	"fmt"
	"log"
	"marketflow/internal/api/handlers"
	"marketflow/internal/domain"
	"marketflow/internal/packages/envzilla"
	"marketflow/internal/service"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	checkFlags()

	if err := envzilla.Loader(".env"); err != nil {
		log.Fatalf("Config file load error: %s", err.Error())
	}
}

// Setup function sets connection to the adapters
func Setup(db domain.Database, cacheMemory domain.CacheMemory, datafetchServ *service.DataModeServiceImp) *http.ServeMux {
	modeHandler := handlers.NewSwitchModeHandler(datafetchServ)
	marketHandler := handlers.NewMarketDataHandler(datafetchServ)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /mode/{mode}", modeHandler.SwitchMode) // Switch to MODE

	mux.HandleFunc("GET /health", modeHandler.CheckHealth) // Returns system status

	mux.HandleFunc("GET /prices/{metric}/{symbol}", marketHandler.ProcessMetricQueryByAll)
	mux.HandleFunc("GET /prices/{metric}/{exchange}/{symbol}", marketHandler.ProcessMetricQueryByExchange)
	fmt.Println(time.Now())
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
