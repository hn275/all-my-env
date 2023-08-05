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
