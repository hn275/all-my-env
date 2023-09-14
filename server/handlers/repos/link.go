package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
	"github.com/jmoiron/sqlx"
)

// returns 201 on success, no body
// { "message": "err" } otherwise
func Link(w http.ResponseWriter, r *http.Request) {
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
		log.Println(res.Status)
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

	// SAVE TO DB
	repo.UserID = user.ID
	repo.ID = repoInfo.ID

	db := database.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTxx(ctx, &sql.TxOptions{Isolation: 0, ReadOnly: false})
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	err = createRepo(tx, &repo)
	if err != nil {
		e, ok := err.(*mysql.MySQLError)
		duplicateKeyCode := uint16(1062)
		if !ok || e.Number != duplicateKeyCode {
			api.NewResponse(w).ServerError(err.Error())
			return
		}

		api.NewResponse(w).Status(http.StatusBadRequest).Error("Repository is already linked.")
		return
	}

	err = tx.Commit()
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	api.NewResponse(w).Status(http.StatusCreated).Done()
}

func createRepo(tx *sqlx.Tx, repo *database.Repository) error {
	// write to `repositories`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `INSERT INTO repositories (id,full_name,user_id) VALUES (:id,:full_name,:user_id);`
	_, err := tx.NamedExecContext(ctx, q, repo)
	if err != nil {
		return err
	}

	// write owner id to `permissions`
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q = `INSERT INTO permissions (user_id,repository_id) VALUES (?,?);`
	_, err = tx.ExecContext(ctx, q, repo.UserID, repo.ID)
	if err != nil {
		return err
	}

	return nil
}

func (cx *githubClient) getRepo(repoName, userToken string, repo *Repository) (int, error) {
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#get-a-repository
	res, err := gh.New(userToken).Get("/repos/%s", repoName)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return res.StatusCode, nil
	}

	if err := json.NewDecoder(res.Body).Decode(&repo); err != nil {
		return -1, err
	}
	return res.StatusCode, nil
}
