package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Response struct {
	http.ResponseWriter
	code int
}

func New(w http.ResponseWriter) *Response {
	return &Response{w, 0}
}

func (r *Response) ServerError(err error) {
	fmt.Fprint(os.Stderr, err)
	r.WriteHeader(http.StatusInternalServerError)
}

func (r *Response) Status(c int) *Response {
	r.code = c
	return r
}

func (r *Response) JSON(data interface{}) {
	if r.code == 0 {
		panic("response code not set")
	}

	r.WriteHeader(r.code)
	r.Header().Add("content-type", "application/json")
	if err := json.NewEncoder(r).Encode(&data); err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		os.Stderr.WriteString(err.Error())
	}
}

func (r *Response) Error(errMessage string) {
	if r.code == 0 {
		panic("response code not set")
	}

	r.WriteHeader(r.code)
	r.Header().Add("content-type", "application/json")
	msg := map[string]string{"error": errMessage}
	if err := json.NewEncoder(r).Encode(msg); err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		os.Stderr.WriteString(err.Error())
	}
}

func (r *Response) Text(t string) {
	if r.code == 0 {
		panic("response code not set")
	}

	r.WriteHeader(r.code)
	r.Header().Add("content-type", "application/text")
	r.Write([]byte(t))
}

func (r *Response) Done() {
	if r.code == 0 {
		panic("response code not set")
	}

	r.WriteHeader(r.code)
}
