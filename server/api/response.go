package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Response struct {
	http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{w}
}

func (r *Response) ServerError(err error) {
	fmt.Fprint(os.Stderr, err)
	r.WriteHeader(http.StatusInternalServerError)
}

func (r *Response) Status(c int) *Response {
	r.WriteHeader(c)
	return r
}

func (r *Response) JSON(data interface{}) {
	r.Header().Add("content-type", "application/json")
	if err := json.NewEncoder(r).Encode(&data); err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		os.Stderr.WriteString(err.Error())
	}
}

func (r *Response) Error(m string) {
	r.Header().Add("content-type", "application/json")
	msg := map[string]string{"error": m}
	if err := json.NewEncoder(r).Encode(&msg); err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		os.Stderr.WriteString(err.Error())
	}
}

func (r *Response) Text(t string) {
	r.Header().Add("content-type", "application/text")
	r.Write([]byte(t))
}
