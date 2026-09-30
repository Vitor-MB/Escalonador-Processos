package controllers

import (
	"encoding/json"
	"escalprocess/internal/dto"
	"escalprocess/internal/response"
	"escalprocess/internal/service"
	"net/http"
)

// Simulate recebe uma requisição de simulação, valida o corpo e retorna o resultado do escalonamento.
func Simulate(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestSimulate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	out, err := service.Simulate(req)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, out)
}

// GetAlgorithms retorna a lista de algoritmos suportados pela API.
func GetAlgorithms(w http.ResponseWriter, r *http.Request) {
	out := service.GetAlgorithms()
	response.WriteJSON(w, http.StatusOK, out)
}
