package repos

import (
	"encoding/json"
	"net/http"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
)

func Link(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	// user, err := jsonwebtoken.GetUser(r)
	_, err := api.NewContext(r).User()
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	var repo database.Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		// TODO: validate request json
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
		return
	}

	// CHECK FOR REPO OWNER
	// get repoInfo
	/*
		var repoInfo Repository
		status, err := ghCx.getRepo(repo.FullName, userID, &repoInfo)
		if err != nil {
			api.NewResponse(w).ServerError(err)
			return
		}

		switch status {
		case http.StatusOK:
			break

		case http.StatusForbidden:
			api.NewResponse(w).Status(status).Error("Forbidden")
			return

		case http.StatusNotFound:
			api.NewResponse(w).
				Status(status).
				Error("Repository not found: %s", repo.FullName)
			return

		default:
			fmt.Fprintf(os.Stderr, "GitHub responded with %d\n", status)
			api.NewResponse(w).
				Status(http.StatusBadGateway).Done()
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
		repo.Url = repoInfo.HTMLURL

		err = db.newRepo(&repo)
		if err == nil {
			api.NewResponse(w).Status(http.StatusCreated).Text("%d", repo.ID)
			return
		}

		pgErr, ok := err.(*pgconn.PgError)
		if !ok {
			api.NewResponse(w).ServerError(err)
			return
		}

		switch pgErr.Code {
		case pgerrcode.ForeignKeyViolation:
			api.NewResponse(w).
				Status(http.StatusBadRequest).
				Error("User not found: %s", user.Login)
			return

		case pgerrcode.UniqueViolation:
			api.NewResponse(w).
				Status(http.StatusBadRequest).
				Error("Repository exists in database: %s", repoInfo.FullName)
			return

		default:
			api.NewResponse(w).ServerError(pgErr)
			return
		}
	*/
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
