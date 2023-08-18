package auth

import (
	"net/http"
	"strconv"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/jsonwebtoken"
)

func TokenValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// VALIDATE COOKIE
		refCookie, err := r.Cookie(api.CookieRefTok)
		if err != nil {
			api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
			return
		}
		if err := refCookie.Valid(); err != nil {
			api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		}

		tok := r.Header.Get("Authorization")
		accessToken, err := getToken(tok)
		if err != nil {
			api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
			return
		}

		// DECODE JWT
		token, err := jsonwebtoken.NewDecoder().Decode(accessToken)
		if err != nil {
			api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
			return
		}
		if len(token.Audience) != 1 {
			api.NewResponse(w).Status(http.StatusForbidden).Error("invalid token")
			return
		}

		// GETTING GITHUB TOKEN FROM REQUEST
		userID, err := strconv.ParseUint(token.Subject, 10, 64)
		if err != nil {
			api.NewResponse(w).ServerError(err.Error())
			return
		}
		userToken, err := decodeAccessToken(userID, token.AccessToken)
		if err != nil {
			api.NewResponse(w).Status(http.StatusUnauthorized).Error(err.Error())
			return
		}
		userLogin := token.Audience[0]

		// ATTACH TO REQUEST
		ctx := api.NewContext(r).SetUser(userID, userToken, userLogin)
		next.ServeHTTP(w, ctx.Request)
	})
}
