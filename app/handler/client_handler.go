package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/oliveirabalsa/go-globalhitss-be/app/dto"
	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"github.com/oliveirabalsa/go-globalhitss-be/app/usecase"
)

type ClientHandler struct {
	ClientUsecase usecase.ClientUsecase
}

// @Summary Create a new client
// @Description Creates a new client with the provided data
// @Tags clients
// @Accept json
// @Produce json
// @Param client body model.Client true "Client object to be created"
// @Success 200 {string} string "Your data has been received and is being processed."
// @Failure 400 {string} string "Bad Request"
// @Router /api/v1/clients [post]
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

// @Summary Update an existing client
// @Description Updates an existing client with the provided data
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Param client body model.Client true "Client object to be updated"
// @Success 200 {string} string "Your data has been received and is being processed."
// @Failure 400 {string} string "Bad Request"
// @Router /api/v1/clients/{id} [patch]
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
// @Success 200 {array} dto.PaginationResponse
// @Router /api/v1/clients [get]
func (h *ClientHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 10
	}

	clients, totalPages, err := h.ClientUsecase.GetClients(page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.PaginationResponse{
		Data:       clients,
		Page:       page,
		NextPage:   page + 1,
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Get a client by ID
// @Description Get a single client by ID from the database
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {object} model.Client
// @Failure 404 {string} string "Client not found"
// @Failure 400 {string} string "Bad Request"
// @Router /api/v1/clients/{id} [get]
func (h *ClientHandler) GetClientByID(w http.ResponseWriter, r *http.Request) {
	clientId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	client, err := h.ClientUsecase.GetClientByID(clientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}

// @Summary Delete an existing client
// @Description Deletes an existing client by ID
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {string} string "Your data has been received and is being processed."
// @Failure 400 {string} string "Invalid client ID"
// @Failure 400 {string} string "Bad Request"
// @Router /api/v1/clients/{id} [delete]
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
