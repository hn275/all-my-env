package envhubtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

func AllowedRequestMethods(m ...string) []string {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPatch,
		http.MethodPut,
		http.MethodHead,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
	for _, allowedMethod := range m {
		for i, method := range methods {
			if allowedMethod != method {
				continue
			}
			left := methods[0:i]
			if i == len(methods)-1 {
				methods = left
			} else {
				right := methods[i+1:]
				methods = append(left, right...)
			}
			break
		}
	}
	return methods
}

func RequestWithParam(
	method string,
	path string,
	urlParams map[string]string,
	body io.Reader,
) *http.Request {
	r := httptest.NewRequest(method, path, body)

	rctx := chi.NewRouteContext()
	for k, v := range urlParams {
		rctx.URLParams.Add(k, v)
	}

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
