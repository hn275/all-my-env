package repos

import (
	"encoding/json"
	"log"
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
	if r.Method != http.MethodDelete {
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
	var partialData map[string]interface{}

	// Decode the JSON into a map
	if err := json.NewDecoder(r.Body).Decode(&partialData); err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
		return
	}

	// Populate the specific fields in the 'repo' struct
	if id, ok := partialData["repoID"].(uint64); ok {
		repo.ID = id
	}
	if fullName, ok := partialData["repoName"].(string); ok {
		repo.FullName = fullName
	}

	// Now, 'repo' should contain the decoded values
	// log.Println("Received repo:", repo)

	// CHECK FOR REPO OWNER
	// get repoInfo
	var repoInfo Repository
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#get-a-repository
	res, err := gh.New(user.Token).Get("/repos/%s", repo.FullName)
	defer res.Body.Close()
	if err != nil {
		log.Println("Error getting repository info:", err)
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	if res.StatusCode != http.StatusOK {
		log.Println("Error from response status:", res.Status)
		api.NewResponse(w).ForwardBadRequest(res)
		return
	}
	if json.NewDecoder(res.Body).Decode(&repoInfo); err != nil {
		log.Println("Error decoding repository info:", err)
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

	// Delete the repository and associated permissions
	if err := db.deleteRepo(repoInfo.ID, repoInfo.Owner.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.ForeignKeyViolation:
				api.NewResponse(w).
					Status(http.StatusNotFound).
					Error("Repository not linked")
			default:
				log.Printf("Database error: %v", pgErr.Message)
				api.NewResponse(w).ServerError("A database error occurred")
			}
		} else {
			log.Printf("Unexpected error: %v", err)
			api.NewResponse(w).ServerError("An unexpected error occurred")
		}
		return

	}

	// Respond with success
	api.NewResponse(w).Status(http.StatusNoContent).Done()

}
