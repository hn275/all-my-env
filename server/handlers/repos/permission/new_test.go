package permission

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPermissionDiff(t *testing.T) {
	fromDB := []uint64{2, 3, 1}
	fromRequest := []uint64{2, 4, 3, 5}

	expected := permDiff{
		revoked: []uint64{1},
		granted: []uint64{4, 5},
	}
	d := getPermssionDiff(fromDB, fromRequest)

	for i := range d.revoked {
		assert.Equal(t, expected.revoked[i], d.revoked[i])
	}

	for i := range d.granted {
		assert.Equal(t, expected.granted[i], d.granted[i])
	}
}
