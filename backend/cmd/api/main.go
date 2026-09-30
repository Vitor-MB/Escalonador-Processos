package main

import (
	"escalprocess/internal/router"
	"log"
	"net/http"
	"time"
)

// main inicializa e sobe a API HTTP do escalonador.
func main() {
	// Cria o roteador com os endpoints da aplicação.
	router := router.GenerateRouter()

	// Porta em que a API será publicada.
	port := "8080"

	// Configura o servidor com timeout de leitura do cabeçalho HTTP.
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Registra a inicialização da API e mantém o processo em execução.
	log.Printf("Iniciando a API na porta: %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
