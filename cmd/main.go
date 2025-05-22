package main

import (
	"log"
	"marketflow/internal/app"
	"marketflow/internal/domain"
	"net/http"
)

func main() {
	router := app.Setup()

	log.Printf("Starting server at %s... \n", *domain.Port)
	if err := http.ListenAndServe("localhost:"+*domain.Port, router); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}
