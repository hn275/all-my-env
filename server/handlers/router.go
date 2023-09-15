package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hn275/envhub/server/handlers/auth"
	"github.com/hn275/envhub/server/handlers/repos"
	"github.com/hn275/envhub/server/handlers/repos/permission"
	"github.com/hn275/envhub/server/handlers/repos/variables"
)

func New() *chi.Mux {
	r := chi.NewMux()

	// CONFIG
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowCredentials: true,
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

	// ROUTES
	r.Route("/auth", func(r chi.Router) {
		r.Handle("/", http.HandlerFunc(auth.LogIn))
		r.Handle("/logout", http.HandlerFunc(auth.Logout))
		r.Handle("/refresh", http.HandlerFunc(auth.RefreshToken))
	})

	r.Route("/repos", func(r chi.Router) {
		r.Use(auth.TokenValidator)
		r.Handle("/", http.HandlerFunc(repos.Index))
		r.Handle("/link", http.HandlerFunc(repos.Link))
		r.Route("/{repoID}", func(r chi.Router) {
			r.Handle("/", http.HandlerFunc(variables.Index))
			r.Route("/variables", func(r chi.Router) {
				r.Handle("/new", http.HandlerFunc(variables.NewVariable))
				r.Handle("/{variableID}", http.HandlerFunc(variables.ID))
			})
			r.Route("/permissions", func(r chi.Router) {
				r.Handle("/", http.HandlerFunc(permission.NewPermission))
			})
		})
	})

	return r
}
