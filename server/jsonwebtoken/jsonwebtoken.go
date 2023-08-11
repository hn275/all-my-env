package jsonwebtoken

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/lib"
)

var (
	secret string
	Secret string
	// Decoder JsonWebToken
	decoder JsonWebTokenDecoder
	encoder JsonWebTokenEncoder

	ErrInvalidToken     error = errors.New("token expired")
	ErrInvalidTokenAlgo error = errors.New("invalid signing algorithm")
)

type AuthClaim struct {
	AccessToken string `json:"access_token"`
	*jwt.RegisteredClaims
}

func init() {
	secret = lib.Getenv("JWT_SECRET")
	Secret = lib.Getenv("JWT_SECRET")

	decoder = &Decoder{}
	encoder = &Encoder{}
}

type JsonWebTokenDecoder interface {
	Decode(string) (*AuthClaim, error)
}
type Decoder struct{}

func NewDecoder() JsonWebTokenDecoder {
	return decoder
}

func (d *Decoder) Decode(token string) (*AuthClaim, error) {
	tok, err := jwt.ParseWithClaims(token, &AuthClaim{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, ErrInvalidTokenAlgo
		}
		return []byte(secret), nil
	})

	if !tok.Valid {
		return nil, ErrInvalidToken
	}

	if err != nil {
		return nil, err
	}

	a, ok := tok.Claims.(*AuthClaim)
	if !ok {
		return nil, errors.New("invalid auth type")
	}

	return a, nil
}

type JsonWebTokenEncoder interface {
	Encode(userID uint64, maskedToken, aud string) (string, error)
}
type Encoder struct{}

func NewEncoder() JsonWebTokenEncoder {
	return encoder
}

func (e *Encoder) Encode(userID uint64, maskedToken, aud string) (string, error) {
	c := AuthClaim{
		AccessToken: maskedToken,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "EnvHub",
			Subject:   fmt.Sprintf("%d", userID),
			Audience:  []string{aud},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().UTC().Add(time.Hour * 24 * 7)},
			NotBefore: &jwt.NumericDate{Time: time.Now().UTC()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(secret))
}
