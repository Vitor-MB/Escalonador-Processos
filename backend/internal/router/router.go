package router

import (
	"net/http"
)

// Funcao GenerateRouter cria um roteador HTTP com os endpoints definidos na lista de rotas
func GenerateRouter() http.Handler {
	mux := http.NewServeMux()
	configure(mux)
	return withCORS(mux)
}

// Configura o roteador com os endpoints definidos na lista de rotas
func configure(mux *http.ServeMux) {
	for _, route := range routes {
		mux.HandleFunc(route.method+" "+route.URI, route.function)
	}
}

// Adiciona suporte a CORS para permitir requisições de diferentes origens
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
