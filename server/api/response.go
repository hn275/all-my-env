package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	CookieRefTok string = "refresh_token"
)

type Response struct {
	http.ResponseWriter
	status int
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{w, 0}
}

func (r *Response) SetCookie(c *http.Cookie) *Response {
	http.SetCookie(r.ResponseWriter, c)
	return r
}

func (r *Response) Status(c int) *Response {
	r.status = c
	return r
}

// only call this when there is no response body
func (r *Response) Done() {
	r.WriteHeader(r.status)
}

func (r *Response) Header(k, v string) *Response {
	r.ResponseWriter.Header().Add(k, v)
	return r
}

func (r *Response) JSON(data interface{}) {
	r.ResponseWriter.Header().Add("Content-Type", "application/json")
	r.WriteHeader(r.status)
	if err := json.NewEncoder(r).Encode(&data); err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		os.Stderr.WriteString(err.Error())
	}
}

func (r *Response) Text(t string, a ...any) {
	r.ResponseWriter.Header().Add("content-type", "application/text")
	r.WriteHeader(r.status)
	r.Write([]byte(fmt.Sprintf(t, a...)))
}

// ERROR HANLDING

func (r *Response) Error(m string, a ...any) {
	r.ResponseWriter.Header().Add("content-type", "application/json")
	r.WriteHeader(r.status)
	msg := map[string]string{"message": fmt.Sprintf(m, a...)}
	if err := json.NewEncoder(r).Encode(&msg); err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		os.Stderr.WriteString(err.Error())
	}
}

func (r *Response) ServerError(m string, a ...any) {
	fmt.Fprintf(os.Stderr, "[ERROR] - "+m, a...)
	r.WriteHeader(http.StatusInternalServerError)
}

func (r *Response) ForwardBadRequest(res *http.Response) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(res.Body)
	if err != nil {
		r.ServerError(err.Error())
		return
	}
	r.WriteHeader(res.StatusCode)
	r.Write(buf.Bytes())
}
