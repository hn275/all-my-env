package variables

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)

	case http.MethodDelete:
		handleUnlink(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
