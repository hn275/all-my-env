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

func refreshToken(userID uint64) (string, error) {
	var buf [16]byte
	binary.BigEndian.PutUint64(buf[:8], userID)
	binary.LittleEndian.PutUint64(buf[8:16], uint64(time.Now().UTC().Unix()))

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

func getUint(b []byte) uint64 {
	l := len(b) - 1
	num := uint64(0)
	for i := l; i >= 0; i-- {
		byteOffset := (l - i) * 8
		num |= uint64(b[i]) << byteOffset
	}
	return num
}

func validateAuthToken(t string) (string, error) {
	if t == "" {
		return "", errors.New("missing auth token")
	}

	auth := strings.Split(t, " ")
	if len(auth) != 2 {
		return "", errors.New("illegal auth token")
	}

	if !strings.EqualFold(auth[0], "bearer") {
		return "", errors.New("illegal token type")
	}

	return auth[1], nil
}

func encodeAccessToken(userID uint64, accessToken string) (string, error) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], userID)
	tok, err := crypto.Encrypt(crypto.UserTokenKey, []byte(accessToken), buf[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(tok), nil
}

func decodeAccessToken(userID uint64, maskedToken string) (string, error) {
	accessToken, err := base64.StdEncoding.DecodeString(maskedToken)
	if err != nil {
		return "", err
	}

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], userID)

	tok, err := crypto.Decrypt(crypto.UserTokenKey, []byte(accessToken), buf[:])
	if err != nil {
		return "", err
	}
	return string(tok), nil
}
