package variables

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/db"
	"github.com/hn275/envhub/server/gh"
	jwt "github.com/hn275/envhub/server/jsonwebtoken"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type permission struct {
	allowed bool
	err     error
}

func (d *variableHandler) NewVariable(w http.ResponseWriter, r *http.Request) {
	ghChan := make(chan permission)
	defer close(ghChan)

	// VALIDATE REQUEST
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	user, err := jwt.GetUser(r)
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Done()
		return
	}

	var body EnvVariable
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error(err.Error())
		return
	}

	// CHECKS IF USER IS A COLLABORATOR
	go getRepoAccess(ghChan, body.RepoURL, user)

	// SERIALIZE VARIABLE
	s := chi.URLParam(r, "id")
	repoID, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error(err.Error())
		return
	}

	envVar, err := body.Cipher(uint32(repoID))
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}
	envVar.CreatedAt = db.TimeNow()
	envVar.UpdatedAt = db.TimeNow()

	// WRITE TO DB
	for {
		select {
		case access := <-ghChan:
			if access.err != nil {
				api.NewResponse(w).ServerError(err)
				return
			}

			if !access.allowed {
				api.NewResponse(w).Status(http.StatusForbidden).Done()
				return
			}

			err := d.Create(envVar).Error
			if err == nil {
				api.NewResponse(w).Status(http.StatusCreated).Done()
				return
			}

			pgErr, ok := err.(*pgconn.PgError)
			if !ok {
				api.NewResponse(w).ServerError(err)
				return
			}

			if pgErr.Code == pgerrcode.UniqueViolation {
				api.NewResponse(w).
					Status(http.StatusConflict).
					Error(pgErr.Error())
				return
			}

			api.NewResponse(w).ServerError(pgErr)
			return

		default:
			continue
		}
	}
}

// Endpoint:
// https://docs.github.com/en/rest/collaborators/collaborators?apiVersion=2022-11-28#check-if-a-user-is-a-repository-collaborator
func getRepoAccess(c chan<- permission, repoURL string, u *jwt.GithubUser) {
	if repoURL[0] == '/' {
		repoURL = repoURL[1:]
	}

	g := gh.New(u.Token)
	url := fmt.Sprintf("/repos/%s/collaborators/%s", repoURL, u.Login)
	res, err := g.Get(url)
	if err != nil {
		buf := permission{false, err}
		c <- buf
		return
	}

	fmt.Println("DELETE ME:")
	fmt.Printf("Status [%s]\n", res.Status)

	switch res.StatusCode {
	case http.StatusNoContent:
		c <- permission{true, nil}
		return

	case http.StatusNotFound:
		c <- permission{false, nil}
		return

	default:
		c <- permission{false, fmt.Errorf("Github API reponses [%v]", res.Status)}
		return
	}
}
