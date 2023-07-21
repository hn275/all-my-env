package variables

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenID(t *testing.T) {
	go RefreshVariableCounter()

	iteration := math.MaxUint16
	repoID := uint32(1)
	idSet := make(map[string]bool)
	for i := 0; i < iteration; i++ {
		idSet[genVariableID(repoID)] = true
	}
	assert.Equal(t, 1, len(counterMap))

	time.Sleep(time.Second)
	assert.Equal(t, 0, len(counterMap))

	for i := 0; i < iteration; i++ {
		idSet[genVariableID(repoID)] = true
	}
	assert.Equal(t, 1, len(counterMap))
	assert.Equal(t, iteration*2, len(idSet))
}
