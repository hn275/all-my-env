package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/hn275/envhub/server/crypto"
)

func refreshToken(userID uint32) (string, error) {
	var buf [16]byte
	binary.BigEndian.PutUint32(buf[:4], userID)
	binary.LittleEndian.PutUint64(buf[4:12], uint64(time.Now().UTC().Unix()))

	// truncating leading - trailing 0s
	l := 0
	for {
		if buf[l] != 0 {
			break
		}
		l++
	}
	r := len(buf) - 1
	for {
		if buf[r] != 0 {
			break
		}
		r--
	}
	for i := 0; i < r; i++ {
		a := i + l
		buf[i] = buf[a]
	}

	// random bytes filler
	_, err := io.ReadFull(rand.Reader, buf[r-l+1:])
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(buf[:]), nil
}

func encodeAccessToken(userID uint32, accessToken string) (string, error) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], userID)
	tok, err := crypto.Encrypt(crypto.UserTokenKey, []byte(accessToken), buf[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(tok), nil
}

func decodeAccessToken(userID uint32, maskedToken string) (string, error) {
	accessToken, err := base64.StdEncoding.DecodeString(maskedToken)
	if err != nil {
		return "", err
	}

	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], userID)

	tok, err := crypto.Decrypt(crypto.UserTokenKey, []byte(accessToken), buf[:])
	if err != nil {
		return "", err
	}
	return string(tok), nil
}

func getToken(token string) (string, error) {
	t := strings.Split(token, " ")
	if len(t) != 2 {
		return "", errors.New("invalid auth token")
	}

	if !strings.EqualFold(t[0], "bearer") {
		return "", errors.New("illegal token type")
	}

	return t[1], nil
}
