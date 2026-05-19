package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	app, err := InitializeService()
	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf(":%d", app.Cfg.Port)
	log.Fatal(http.ListenAndServe(addr, app.Engine))
}
