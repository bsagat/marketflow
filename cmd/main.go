package main

import (
	"context"
	"log"
	"log/slog"
	cache "marketflow/internal/adapters/cacheMemory"
	datafetcher "marketflow/internal/adapters/dataFetcher"
	"marketflow/internal/adapters/repository"
	"marketflow/internal/app"
	"marketflow/internal/domain"
	"marketflow/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	srv, cleanup := setupApp()
	defer cleanup()

	startServer(srv)

	waitForShutdown()

	shutdownServer(srv)

	slog.Info("App is closed...")
}

func setupApp() (*http.Server, func()) {
	cacheMemory := cache.ConnectCacheMemory()
	repo := repository.ConnectDB()
	datafetch := datafetcher.NewLiveModeFetcher()
	datafetchServ := service.NewDataFetcher(datafetch, repo, cacheMemory)

	if err := datafetchServ.ListenAndSave(); err != nil {
		slog.Error("Failed to start data fetcher", "error", err)
		datafetch.Close()
		os.Exit(1)
	}

	router := app.Setup(repo, cacheMemory, datafetchServ)
	srv := &http.Server{
		Addr:    "localhost:" + *domain.Port,
		Handler: router,
	}

	cleanup := func() {
		slog.Info("Cleaning up resources...")
		datafetchServ.StopListening()
		cacheMemory.Cache.Close()
		repo.Db.Close()
	}

	return srv, cleanup
}

func startServer(srv *http.Server) {
	go func() {
		slog.Info("Starting server at " + *domain.Port + "...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: ", err.Error())
		}
	}()
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	slog.Info("Shutdown signal received...")
}

func shutdownServer(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("Shutting down HTTP server...")
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	} else {
		slog.Info("Server gracefully stopped.")
	}
}
