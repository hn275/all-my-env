package repos

import (
	"encoding/json"
	"net/http"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// returns 201 on success, no body
// { "message": "err" } otherwise
func Unlink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	user, err := api.NewContext(r).User()
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	// only `repo.ID` and `repo.FullName` are sent
	var repo database.Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
		return
	}

	// CHECK FOR REPO OWNER
	// get repoInfo
	var repoInfo Repository
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#get-a-repository
	res, err := gh.New(user.Token).Get("/repos/%s", repo.FullName)
	defer res.Body.Close()
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	if res.StatusCode != http.StatusOK {
		api.NewResponse(w).ForwardBadRequest(res)
		return
	}
	if json.NewDecoder(res.Body).Decode(&repoInfo); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	isOwner := user.ID == repoInfo.Owner.ID
	if !isOwner {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("Not the repository's owner")
		return
	}

	// DELETE FROM DB
	err = db.deleteRepo(repo.ID)
	if err == nil {
		api.NewResponse(w).Status(http.StatusNoContent).Done()
		return
	}

	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	switch pgErr.Code {
	case pgerrcode.ForeignKeyViolation:
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("Repository not linked")
	default:
		api.NewResponse(w).ServerError(err.Error())
	}

}
