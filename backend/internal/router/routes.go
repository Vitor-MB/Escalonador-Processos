package router

import (
	"escalprocess/internal/controllers"
	"net/http"
)

// Struct Rota para armazenar informações sobre cada endpoint
type route struct {
	URI      string
	method   string
	function func(http.ResponseWriter, *http.Request)
}

// Lista de endpoints disponíveis na aplicação
var routes = []route{
	{
		URI:      "/simulate",
		method:   "POST",
		function: controllers.Simulate,
	},
	{
		URI:      "/algorithm",
		method:   "GET",
		function: controllers.GetAlgorithms,
	},
}
