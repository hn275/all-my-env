package api

import (
	"context"
	"errors"
	"net/http"
)

const (
	CtxUser      string = "user"
	CtxUserID    string = "user_id"
	CtxUserToken string = "user_token"

	appCtx string = "app_ctx"
)

type RequestCtx struct {
	*http.Request
}

type UserContext struct {
	ID    uint64
	Token string
}

func NewContext(r *http.Request) *RequestCtx {
	return &RequestCtx{r}
}

func (r *RequestCtx) SetUser(userID uint64, userToken string) *RequestCtx {
	ctx := context.WithValue(r.Context(), appCtx, &UserContext{userID, userToken})
	return &RequestCtx{r.WithContext(ctx)}
}

func (r *RequestCtx) User() (*UserContext, error) {
	u, ok := r.Context().Value(appCtx).(*UserContext)
	if !ok {
		return nil, errors.New("user context not found")
	}
	return u, nil
}
