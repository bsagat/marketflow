package main

import (
	"fmt"
	"log"
	datafetcher "marketflow/internal/adapters/dataFetcher"
	"marketflow/internal/app"
)

func main() {
	err := app.SetConfig()
	if err != nil {
		log.Fatal("Config set error: ", err.Error())
	}

	mode := datafetcher.LiveMode{}
	aggregated := mode.SetupDataFetcher()
	if aggregated != nil {
		for data := range aggregated {
			fmt.Println(data)
		}
	}
	select {}
}
