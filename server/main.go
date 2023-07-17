package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hn275/envhub/server/handlers/auth"
	"github.com/hn275/envhub/server/handlers/repos"
)

var port string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func main() {
	r := chi.NewMux()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"}, // TODO: configure this cors
	}))

	r.Route("/auth", func(r chi.Router) {
		r.Handle("/github", http.HandlerFunc(auth.Handler.VerifyToken))
	})

	r.Route("/repos", func(r chi.Router) {
		r.Handle("/all", http.HandlerFunc(repos.Handlers.All))
	})

	log.Println("Listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
