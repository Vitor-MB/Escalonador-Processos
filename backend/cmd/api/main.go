package main

import (
	"escalprocess/internal/router"
	"log"
	"net/http"
	"time"
)

func main() {
	router := router.GenerateRouter()

	port := "8080"

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Iniciando a API na porta: %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
