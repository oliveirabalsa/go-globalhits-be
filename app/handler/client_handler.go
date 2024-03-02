package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"github.com/oliveirabalsa/go-globalhitss-be/app/usecase"
)

type ClientHandler struct {
	ClientUsecase usecase.ClientUsecase
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var client model.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client.ID = uuid.New()

	message, err := h.ClientUsecase.CreateClient(&client)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *ClientHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.ClientUsecase.GetClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}
