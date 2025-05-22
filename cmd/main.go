package main

import (
	"log"
	"marketflow/internal/app"
)

func main() {
	err := app.SetConfig()
	if err != nil {
		log.Fatal("Config set error: ", err.Error())
	}

}
