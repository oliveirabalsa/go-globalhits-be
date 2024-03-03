package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/oliveirabalsa/go-globalhitss-be/app/handler"
	"github.com/oliveirabalsa/go-globalhitss-be/app/repository"
	"github.com/oliveirabalsa/go-globalhitss-be/app/usecase"
	"github.com/oliveirabalsa/go-globalhitss-be/config"
	_ "github.com/oliveirabalsa/go-globalhitss-be/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")

	ch, conn, db := config.InitServices()
	defer conn.Close()
	defer ch.Close()

	clientRepo := repository.NewClientRepository(db)
	clientUsecase := usecase.NewClientUseCase(*clientRepo, ch, os.Getenv("RABBITMQ_QUEUE"))
	clientHandler := handler.ClientHandler{ClientUsecase: *clientUsecase}

	r := chi.NewRouter()
	r.Mount("/swagger/", httpSwagger.WrapHandler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/clients", clientHandler.CreateClient)
		r.Get("/clients", clientHandler.GetClients)
		r.Get("/clients/{id}", clientHandler.GetClientByID)
		r.Patch("/clients/{id}", clientHandler.UpdateClient)
		r.Delete("/clients/{id}", clientHandler.DeleteClient)
	})

	log.Println("Starting server on :" + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
