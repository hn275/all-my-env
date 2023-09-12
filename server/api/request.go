package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	appCtx string = "app_ctx"
)

type RequestCtx struct {
	*http.Request
}

type UserContext struct {
	ID    uint32
	Token string
	Login string
}

func NewContext(r *http.Request) *RequestCtx {
	return &RequestCtx{r}
}

func (r *RequestCtx) SetUser(userID uint32, userToken, userLogin string) *RequestCtx {
	u := &UserContext{userID, userToken, userLogin}
	ctx := context.WithValue(r.Context(), appCtx, u)
	return &RequestCtx{r.WithContext(ctx)}
}

func (r *RequestCtx) User() (*UserContext, error) {
	u, ok := r.Context().Value(appCtx).(*UserContext)
	if !ok {
		return nil, errors.New("user context not found")
	}
	return u, nil
}

func (r *RequestCtx) Query(k string) (string, error) {
	q := r.Request.URL.Query().Get(k)
	if q == "" {
		return "", fmt.Errorf("query %s not found", k)
	}
	return q, nil
}
