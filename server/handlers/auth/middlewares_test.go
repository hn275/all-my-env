package auth

import (
	"encoding/binary"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetUInt(t *testing.T) {
	testCases := []uint64{
		420,
		69,
		123,
		3_546,
		123,
		345,
		56_456,
		113_546_565,
	}

	for _, c := range testCases {
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], c)
		assert.Equal(t, c, getUint(buf[:]))
	}
}

func TestTokenChecker(t *testing.T) {
	r := chi.NewMux()
	r.Use(TokenValidator)
	t.Error("implement this")
}
