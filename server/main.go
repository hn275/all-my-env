package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hn275/envhub/server/handlers/auth"
)

var port string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func main() {
	mux := chi.NewMux()
	mux.Use(middleware.Logger)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"}, // TODO: configure this cors
	}))

	mux.Route("/auth", auth.Router)

	log.Println("Listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
