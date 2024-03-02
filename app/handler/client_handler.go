package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"github.com/oliveirabalsa/go-globalhitss-be/app/usecase"
	"github.com/oliveirabalsa/go-globalhitss-be/queue"
	"github.com/streadway/amqp"
)

type ClientHandler struct {
	ClientUsecase usecase.ClientUsecase
	QueueChannel  *amqp.Channel
	QueueName     string
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var client model.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client.ID = uuid.New()
	clientData, err := json.Marshal(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	message := &queue.Message{
		Action: "create_client",
		Data:   clientData,
	}
	err = queue.PublishMessage(h.QueueChannel, h.QueueName, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}
