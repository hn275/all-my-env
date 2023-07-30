package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hn275/envhub/server/database"
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
	routerConfig(r)

	// refresh variable counter every second
	go database.RefreshVariableCounter()

	/* ROUTERS */
	r.Route("/auth", func(r chi.Router) {
		r.Handle("/github", http.HandlerFunc(auth.Handler.VerifyToken))
	})

	r.Route("/repo", func(r chi.Router) {
		r.Handle("/", http.HandlerFunc(repos.Index))
		r.Handle("/link", http.HandlerFunc(repos.Link))

		r.Route("/{repoID}", func(r chi.Router) {
			r.Route("/variables", func(r chi.Router) {
				r.Handle("/", http.HandlerFunc(variables.Handlers.Index))
				r.Handle("/new", http.HandlerFunc(variables.Handlers.NewVariable))
			})
		})
	})

	log.Println("Listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func routerConfig(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
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
}
