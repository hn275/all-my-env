package crypto

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	plaintext := "hello, world"
	variableID := uint32(1)

	ciphertext, err := Encrypt(plaintext, variableID)
	assert.Nil(t, err)
	assert.NotEqual(t, string(ciphertext), plaintext)

	decrypted, err := Decrypt(ciphertext, variableID)
	assert.Nil(t, err)
	assert.Equal(t, string(decrypted), plaintext)

	_, err = Decrypt(ciphertext, variableID+1)
	assert.NotNil(t, err)

	_, err = Decrypt(ciphertext, variableID)
	assert.Nil(t, err)

	_, err = Decrypt(ciphertext[:31], variableID)
	assert.NotNil(t, err)
}

func FuzzDecrypt(f *testing.F) {
	testcases := [][]byte{
		[]byte("hello, world"),
		[]byte(""),
		[]byte(""),
		[]byte(""),
		[]byte("hello, world"),
	}

	for _, testcase := range testcases {
		f.Add(testcase)
	}

	f.Fuzz(func(t *testing.T, ciphertext []byte) {
		_, err := Decrypt(ciphertext, 1)
		if err != nil {
			t.Skip()
		}
	})
}

func FuzzEncrypt(f *testing.F) {
	testcases := []struct {
		value      string
		variableID uint32
	}{
		{
			value:      "hello, world",
			variableID: 1,
		},
		{
			value:      "",
			variableID: 0,
		},
		{
			value:      "",
			variableID: 0,
		},
		{
			value:      "",
			variableID: 0,
		},
		{
			value:      "hello, world",
			variableID: math.MaxInt32,
		},
	}

	for _, testcase := range testcases {
		f.Add(testcase.value, testcase.variableID)
	}

	f.Fuzz(func(t *testing.T, value string, variableID uint32) {
		_, err := Encrypt(value, variableID)
		if err != nil {
			t.Skip()
		}
	})
}

func FuzzEncryptDecrypt(f *testing.F) {
	testcases := []struct {
		value      string
		variableID uint32
	}{
		{
			value:      "hello, world",
			variableID: 1,
		},
		{
			value:      "",
			variableID: 0,
		},

		{
			value:      "hello, world",
			variableID: math.MaxInt32,
		},
	}

	for _, testcase := range testcases {
		f.Add(testcase.value, testcase.variableID)
	}

	f.Fuzz(func(t *testing.T, value string, variableID uint32) {
		ciphertext, err := Encrypt(value, variableID)
		if err != nil {
			t.Skip()
		}
		_, err = Decrypt(ciphertext, variableID)
		if err != nil {
			t.Skip()
		}
	})
}
