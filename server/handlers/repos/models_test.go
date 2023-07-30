// go build:testing
package repos

import (
	"github.com/hn275/envhub/server/database"
)

type mockRepoDb struct {
	findRepoErr error
	newRepoErr  error
}

func (repoDB *mockRepoDb) newRepo(_ *database.Repository) error {
	if repoDB.newRepoErr != nil {
		return repoDB.newRepoErr.(error)
	}
	return nil
}

func (db *mockRepoDb) findRepo(userID uint64, ids []uint64) ([]uint64, error) {
	if db.findRepoErr != nil {
		return nil, db.findRepoErr
	}
	repoIDs := []uint64{1}
	return repoIDs, nil
}
