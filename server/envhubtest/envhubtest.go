package envhubtest

import (
	"net/http"
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
