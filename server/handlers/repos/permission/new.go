package permission

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

// request json contains all the user ids that should have write access:
// { "userIDs": []uint64 }
func NewPermission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	user, err := api.NewContext(r).User()
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	type NewAccessRequest struct {
		UserID uint64
		RepoID uint64
	}

	var newPerm struct {
		UserIDs []uint32 `json:"userIDs"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newPerm); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	// filter out repo owner id

	// CHECK TO FOR REPOSITORY OWNER.
	// query db for github information
	s := chi.URLParam(r, "repoID")
	repoID, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
	}

	var repo struct {
		UserID uint32
	}
	err = db.Model(&database.Repository{}).
		Select([]string{"user_id"}).
		Where("id = ?", repoID).
		First(&repo).
		Error

	switch err {
	case nil:
		break

	case gorm.ErrRecordNotFound:
		api.NewResponse(w).Status(http.StatusNotFound).Error("repository not found")
		return

	default:
		api.NewResponse(w).Status(http.StatusInternalServerError).Error(err.Error())
		return
	}

	if repo.UserID != user.ID {
		api.NewResponse(w).Status(http.StatusForbidden).Error("not repository owner.")
		return
	}

	// WRITE NEW ACCESS ENTRY TO TABLE
	// get all user id with permission
	var p []struct {
		UserID uint32
	}
	err = db.Model(&database.Permission{}).Where("id = ?", repoID).Find(&p).Error
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	uids := make([]uint32, len(p))
	for i := range p {
		uids[i] = p[i].UserID
	}

	// get diff then write to db
	err = getPermssionDiff(repo.UserID, uids, newPerm.UserIDs).updatePermissions(uint32(repoID))
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	api.NewResponse(w).Status(http.StatusOK).Done()
}

type permDiff struct {
	revoked []uint32
	granted []uint32
}

func getPermssionDiff(ownerID uint32, dbPerm, reqPerm []uint32) *permDiff {
	// filter owner id
	i := 0
	for ; i < len(dbPerm); i++ {
		if dbPerm[i] == ownerID {
			break
		}
		i++
	}
	if i < len(dbPerm)-1 && i > 0 {
		dbPerm = append(dbPerm[:i], dbPerm[i+1:]...)
	}

	i = 0
	for ; i < len(reqPerm); i++ {
		if reqPerm[i] == ownerID {
			break
		}
		i++
	}
	if i < len(reqPerm)-1 && i > 0 {
		reqPerm = append(reqPerm[:i], reqPerm[i+1:]...)
	}

	// sort
	sort.SliceStable(dbPerm, func(i, j int) bool {
		return dbPerm[i] < dbPerm[j]
	})

	sort.SliceStable(reqPerm, func(i, j int) bool {
		return reqPerm[i] < reqPerm[j]
	})

	// get lower and upperbound length
	h, l := len(dbPerm), len(reqPerm)
	if h < l {
		l, h = h, l
	}

	// get diff
	p := permDiff{
		revoked: make([]uint32, 0, h),
		granted: make([]uint32, 0, h),
	}

	dbI, reqI := 0, 0
	for dbI < l && reqI < l {
		dbP, reqP := dbPerm[dbI], reqPerm[reqI]
		if dbP == reqP {
			dbI++
			reqI++
		} else if dbP < reqP {
			p.revoked = append(p.revoked, dbP)
			dbI++
		} else {
			p.granted = append(p.granted, reqP)
			reqI++
		}

	}

	// the remainder of db permission is revoked
	for ; dbI < len(dbPerm); dbI++ {
		p.revoked = append(p.revoked, dbPerm[dbI])
	}

	// the remainder of req permission is granted
	for ; reqI < len(reqPerm); reqI++ {
		p.granted = append(p.granted, reqPerm[reqI])
	}

	return &p
}

func (wa *permDiff) updatePermissions(repoID uint32) error {
	return db.Transaction(func(tx *gorm.DB) error {
		type Permision struct {
			RepositoryID uint32
			UserID       uint32
		}

		p := make([]database.Permission, len(wa.revoked))
		for i := range wa.revoked {
			p[i] = database.Permission{
				RepositoryID: repoID,
				UserID:       wa.granted[i],
			}
		}

		err := tx.Where("repository_id = ?", repoID).Delete(&p).Error
		if err != nil {
			return err
		}

		p = make([]database.Permission, len(wa.granted))
		for i := range wa.granted {
			p[i] = database.Permission{
				RepositoryID: repoID,
				UserID:       wa.granted[i],
			}
		}

		return tx.Create(&p).Error
	})
}
