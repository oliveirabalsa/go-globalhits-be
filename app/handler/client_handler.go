package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	if err := client.Validate(); err != nil {
		errorMessage := map[string]string{"error": err.Error()}
		jsonErrorResponse, _ := json.Marshal(errorMessage)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonErrorResponse)
		return
	}

	client.ID = uuid.New()

	message, err := h.ClientUsecase.CreateClient(&client)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	var client model.Client

	clientId, err := uuid.Parse(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client.ID = uuid.New()

	message, err := h.ClientUsecase.UpdateClient(clientId, &client)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// @Summary Get a list of clients
// @Description Get a list of clients from the database
// @Tags clients
// @Accept json
// @Produce json
// @Success 200 {array} model.Client
// @Router /clients [get]
func (h *ClientHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.ClientUsecase.GetClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	clientId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid clientId", http.StatusBadRequest)
		return
	}
	message, err := h.ClientUsecase.DeleteClient(clientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}
