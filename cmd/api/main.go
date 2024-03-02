package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oliveirabalsa/go-globalhitss-be/app/handler"
	"github.com/oliveirabalsa/go-globalhitss-be/app/repository"
	"github.com/oliveirabalsa/go-globalhitss-be/app/usecase"
	"github.com/oliveirabalsa/go-globalhitss-be/db"
	"github.com/oliveirabalsa/go-globalhitss-be/queue"
)

func main() {
	// Initialize database
	db, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize RabbitMQ
	conn, ch, err := queue.NewRabbitMQ()
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	// Initialize repository
	clientRepo := repository.ClientRepository{DB: db}

	// Initialize usecase
	clientUsecase := usecase.ClientUsecase{ClientRepo: clientRepo}

	// Initialize handler
	clientHandler := handler.ClientHandler{
		ClientUsecase: clientUsecase,
		QueueChannel:  ch,
		QueueName:     "globalhitss",
	}

	// Initialize router
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	r.Post("/clients", clientHandler.CreateClient)
	// Add other routes as needed

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
