package database

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenID(t *testing.T) {
	t.Parallel()
	go RefreshVariableCounter()

	iteration := math.MaxUint16
	idSet := make(map[string]bool)
	v := Variable{RepositoryID: 1}

	for i := 0; i < iteration; i++ {
		err := v.GenID()
		assert.Nil(t, err)
		idSet[v.ID] = true
	}
	assert.Equal(t, 1, len(idCounterMap))

	time.Sleep(time.Second)
	assert.Equal(t, 0, len(idCounterMap))

	for i := 0; i < iteration; i++ {
		err := v.GenID()
		assert.Nil(t, err)
		idSet[v.ID] = true
	}
	assert.Equal(t, 1, len(idCounterMap))
	assert.Equal(t, iteration*2, len(idSet))
}

func TestVariableGenID(t *testing.T) {
	v := Variable{Key: "foo", Value: "bar"}
	err := v.GenID()
	assert.Equal(t, ErrRepoMissingRepoID, err)

	v.RepositoryID = 123
	err = v.GenID()
	assert.Nil(t, err)
	assert.NotEmpty(t, v.ID)
	assert.Equal(t, 32, len(v.ID))
}

func TestVariableEncryptDecrypt(t *testing.T) {
	var err error
	v := Variable{Key: "foo", RepositoryID: 123}
	err = v.EncryptValue()
	assert.Equal(t, ErrIDNotGenerated, err)

	err = v.GenID()
	assert.Nil(t, err)

	err = v.EncryptValue()
	assert.Equal(t, ErrValueNotFound, err)

	v.Value = "bar"
	err = v.EncryptValue()
	assert.Nil(t, err)
	assert.NotEqual(t, "bar", v.Value)

	err = v.DecryptValue()
	assert.Nil(t, err)
	assert.Equal(t, "bar", v.Value)
}
