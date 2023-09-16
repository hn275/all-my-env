package variables

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	rCtx := api.NewContext(r)
	user, err := rCtx.User()
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error(err.Error())
		return
	}

	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 32)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("Invalid repository id")
		return
	}

	// GET VARIABLES
	db := database.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get repository
	repo := database.Repository{ID: uint32(repoID)}
	q := `SELECT full_name,user_id FROM repositories WHERE id = ?`
	err = db.QueryRowxContext(ctx, q, repoID).StructScan(&repo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.NewResponse(w).Status(http.StatusNotFound).Error("Repository not found")
			return
		}
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	// get variables
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	variables := make([]database.Variable, 0)
	q = `
	SELECT
		id,
		created_at,
		updated_at,
		variable_key,
		variable_value
	FROM variables
	WHERE repository_id = ?;
	`

	rows, err := db.QueryxContext(ctx, q, repo.ID)
	defer rows.Close()

	for rows.Next() {
		var v database.Variable
		if err := rows.StructScan(&v); err != nil {
			api.NewResponse(w).ServerError(err.Error())
			return
		}
		if err := v.DecryptValue(); err != nil {
			api.NewResponse(w).ServerError(err.Error())
			return
		}
		variables = append(variables, v)
	}

	// get contributor
	contributors, err := getContributors(user, repo.FullName)
	if err != nil {
		api.NewResponse(w).Status(http.StatusBadGateway).Error(err.Error())
		return
	}

	// check for access
	isOwner := repo.UserID == user.ID
	if !isOwner {
		accessOK := false
		for i := range contributors {
			if contributors[i].ID == uint32(user.ID) {
				accessOK = true
			}
		}

		if !accessOK {
			api.NewResponse(w).Status(http.StatusForbidden).Error("Not a contributors.")
			return
		}
	}

	// get contributors with write access
	writeAccess := make(map[uint32]bool)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q = `SELECT user_id FROM permissions WHERE repository_id = ? ORDER BY user_id`
	rows, err = db.QueryxContext(ctx, q, repo.ID)
	defer rows.Close()

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	for rows.Next() {
		var u uint32
		if err := rows.Scan(&u); err != nil {
			api.NewResponse(w).ServerError(err.Error())
			return
		}
		writeAccess[u] = true
	}

	// combine data from gh api and database
	for i, u := range contributors {
		_, ok := writeAccess[uint32(u.ID)]
		contributors[i].WriteAccess = ok
	}

	response := struct {
		Variables    []database.Variable `json:"variables"`
		WriteAccess  bool                `json:"write_access"`
		OwnerID      uint32              `json:"owner_id"`
		IsOwner      bool                `json:"is_owner"`
		Contributors []gh.Contributor    `json:"contributors"`
	}{
		Variables:    variables,
		WriteAccess:  isOwner || writeAccess[user.ID],
		IsOwner:      isOwner,
		OwnerID:      repo.UserID,
		Contributors: contributors,
	}

	api.NewResponse(w).
		//Header("Cache-Control", "max-age=20").
		Status(http.StatusOK).
		JSON(&response)
}

func getContributors(user *api.UserContext, repo string) ([]gh.Contributor, error) {
	res, err := gh.New(user.Token).Get("/repos/%s/collaborators", repo)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var buf bytes.Buffer
		buf.ReadFrom(res.Body)
		fmt.Fprintf(os.Stderr, "GitHub API response %s - %v\n", res.Status, buf.String())
		return nil, errors.New("Github API failed.")
	}

	var c []gh.Contributor
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		return nil, err
	}

	return c, nil
}
