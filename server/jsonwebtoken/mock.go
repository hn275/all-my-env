package jsonwebtoken

import (
	"encoding/base64"
	"encoding/binary"

	"github.com/golang-jwt/jwt/v5"
)

type jwtDecoderMock struct{}

func Mock() {
	// decoder = &jwtDecoderMock{}
}

func (j *jwtDecoderMock) Decode(_ string) (*jwt.Token, error) {
	userID := uint64(1)
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], userID)

	u := &jwt.Token{
		Raw:    "",
		Method: nil,
		Header: map[string]interface{}{},
		Claims: AuthClaim{
			AccessToken: base64.StdEncoding.EncodeToString(buf[:]),
		},
		Signature: []byte{},
		Valid:     false,
	}
	return u, nil
}
