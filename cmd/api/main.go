package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oliveirabalsa/go-globalhitss-be/app/handler"
	"github.com/oliveirabalsa/go-globalhitss-be/app/repository"
	"github.com/oliveirabalsa/go-globalhitss-be/app/usecase"
	"github.com/oliveirabalsa/go-globalhitss-be/config"
)

func main() {
	ch, conn, db := config.InitServices()
	defer conn.Close()
	defer ch.Close()

	clientRepo := repository.NewClientRepository(db)

	clientUsecase := usecase.NewClientUseCase(*clientRepo, ch, "globalhitss")

	clientHandler := handler.ClientHandler{
		ClientUsecase: *clientUsecase,
	}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	r.Post("/clients", clientHandler.CreateClient)
	r.Get("/clients", clientHandler.GetClients)
	// Add other routes as needed

	log.Println("Starting server on :8082...")
	if err := http.ListenAndServe(":8082", r); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
