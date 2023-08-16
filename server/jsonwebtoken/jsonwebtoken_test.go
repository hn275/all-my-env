package jsonwebtoken

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	uid := uint64(123)
	maskedTok := "sometokenhere"
	aud := "foo"

	token, err := NewEncoder().Encode(uid, maskedTok, aud)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	clm, err := NewDecoder().Decode(token)
	assert.Nil(t, err)
	assert.Equal(t, aud, clm.Audience[0])

	decodedID, err := strconv.ParseUint(clm.Subject, 10, 64)
	assert.Nil(t, err)
	assert.Equal(t, uid, decodedID)
}
