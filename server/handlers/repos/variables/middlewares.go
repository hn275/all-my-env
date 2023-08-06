package variables

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/jsonwebtoken"
	"gorm.io/gorm"
)

type RepoContext struct {
	RepoID uint64
	User   *jsonwebtoken.GithubUser
}

func WriteAccessChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 64)
		if err != nil {
			api.NewResponse(w).
				Status(http.StatusBadRequest).
				Error("invalid repository id")
			return
		}

		user, err := jsonwebtoken.GetUser(r)
		if err != nil {
			api.NewResponse(w).Status(http.StatusForbidden).Done()
			return
		}

		repo := database.Repository{ID: repoID}
		err = db.getRepoAccess(&repo, user.ID)
		switch err {
		case nil:
			repoCtx := &RepoContext{repoID, user}
			ctx := context.WithValue(r.Context(), "repoCtx", repoCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
			return

		case gorm.ErrRecordNotFound:
			api.NewResponse(w).
				Status(http.StatusBadRequest).
				Error("Write-access not granted. Please contact repo owner.")
			return

		default:
			api.NewResponse(w).ServerError(err)
			return
		}
	})
}
