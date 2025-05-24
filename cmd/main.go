package main

import (
	"log"
	cache "marketflow/internal/adapters/cacheMemory"
	datafetcher "marketflow/internal/adapters/dataFetcher"
	"marketflow/internal/adapters/repository"
	"marketflow/internal/app"
	"marketflow/internal/domain"
	"net/http"
)

func main() {
	cacheMemory := cache.ConnectCacheMemory()
	repo := repository.ConnectDB()
	datafetch := datafetcher.NewLiveModeFetcher()

	defer cacheMemory.Cache.Close()
	defer repo.Db.Close()
	defer datafetch.Close()

	router := app.Setup(repo, cacheMemory, datafetch)

	log.Printf("Starting server at %s... \n", *domain.Port)
	if err := http.ListenAndServe("localhost:"+*domain.Port, router); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}
