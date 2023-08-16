package jsonwebtoken

import (
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtDecoderMock struct{}

func MockDecoder(m JsonWebTokenDecoder) {
	decoder = m
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
			RegisteredClaims: &jwt.RegisteredClaims{
				Issuer:    "EnvHub",
				Subject:   "1",
				Audience:  []string{"foo"},
				ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(100 * time.Hour)},
				ID:        "",
			},
		},
		Signature: []byte{},
		Valid:     true,
	}
	return u, nil
}
