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
	"github.com/hn275/envhub/server/handlers/repos/variables"
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
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3000/",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3000/",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{
			"Authorization",
			"Accept",
			"Content-Type",
		},
	}))

	r.Route("/auth", func(r chi.Router) {
		r.Handle("/github", http.HandlerFunc(auth.Handler.VerifyToken))
	})

	// refresh variable counter every second
	go variables.RefreshVariableCounter()
	r.Route("/repos", func(r chi.Router) {
		r.Handle("/", http.HandlerFunc(repos.Handlers.All))
		r.Route("/{id}", func(r chi.Router) {
			r.Route("/variables", func(r chi.Router) {
				r.Handle("/new", http.HandlerFunc(variables.Handlers.NewVariable))
			})
		})
	})

	log.Println("Listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
