package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router(r chi.Router) {
	r.Handle("/github", http.HandlerFunc(verify))
}
