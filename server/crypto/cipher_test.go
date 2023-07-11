package crypto

import (
	"math"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	encryptionKey := []byte("01234567890123456789012345678901")
	plaintext := "hello, world"
	variableID := 1

	ciphertext, err := Encrypt(encryptionKey, plaintext, variableID)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	decrypted, err := Decrypt(encryptionKey, ciphertext, variableID)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Fatalf("decrypted text is not same as plaintext: %s", string(decrypted))
	}

	_, err = Decrypt(encryptionKey, ciphertext, variableID+1)
	if err == nil {
		t.Fatalf("decrypted text with wrong variableID")
	}

	_, err = Decrypt(encryptionKey[:31], ciphertext, variableID)
	if err == nil {
		t.Fatalf("decrypted text with wrong key")
	}

	_, err = Decrypt(encryptionKey, ciphertext[:31], variableID)
	if err == nil {
		t.Fatalf("decrypted text with wrong ciphertext")
	}

	_, err = Decrypt(encryptionKey, ciphertext, variableID)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}
}

func FuzzDecrypt(f *testing.F) {
	testcases := []struct {
		encryptionKey []byte
		ciphertext    []byte
	}{
		{
			encryptionKey: []byte("01234567890123456789012345678901"),
			ciphertext:    []byte("hello, world"),
		},
		{
			encryptionKey: nil,
			ciphertext:    []byte(""),
		},
		{
			encryptionKey: []byte{0x00},
			ciphertext:    []byte(""),
		},
		{
			encryptionKey: []byte("01234567890123456789012345678901"),
			ciphertext:    []byte(""),
		},
		{
			encryptionKey: []byte("01234567890123456789012345678901"),
			ciphertext:    []byte("hello, world"),
		},
	}

	for _, testcase := range testcases {
		f.Add(testcase.encryptionKey, testcase.ciphertext)
	}

	f.Fuzz(func(t *testing.T, encryptionKey []byte, ciphertext []byte) {
		_, err := Decrypt(encryptionKey, ciphertext, 1)
		if err != nil {
			t.Skip()
		}
	})
}

func FuzzEncrypt(f *testing.F) {
	testcases := []struct {
		encryptionKey []byte
		value         string
		variableID    int
	}{
		{
			encryptionKey: []byte("01234567890123456789012345678901"),
			value:         "hello, world",
			variableID:    1,
		},
		{
			encryptionKey: nil,
			value:         "",
			variableID:    0,
		},
		{
			encryptionKey: []byte{0x00},
			value:         "",
			variableID:    0,
		},
		{
			encryptionKey: []byte{0x00, 0x00},
			value:         "",
			variableID:    0,
		},
		{
			encryptionKey: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			value:      "\x00",
			variableID: -1,
		},
		{
			encryptionKey: []byte{
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			},
			value:      "\xFF",
			variableID: -1,
		},
		{
			encryptionKey: []byte("01234567890123456789012345678901"),
			value:         "hello, world",
			variableID:    math.MaxInt32,
		},
		{
			encryptionKey: []byte("01234567890123456789012345678901"),
			value:         "hello, world",
			variableID:    math.MinInt32,
		},
	}

	for _, testcase := range testcases {
		f.Add(testcase.encryptionKey, testcase.value, testcase.variableID)
	}

	f.Fuzz(func(t *testing.T, encryptionKey []byte, value string, variableID int) {
		_, err := Encrypt(encryptionKey, value, variableID)
		if err != nil {
			t.Skip()
		}
	})
}

func FuzzEncryptDecrypt(f *testing.F) {
	testcases := []struct {
		encryptionKey []byte
		value         string
		variableID    int
	}{
		{
			encryptionKey: []byte("01234567890123456789012345678901"),
			value:         "hello, world",
			variableID:    1,
		},
		{
			encryptionKey: nil,
			value:         "",
			variableID:    0,
		},
		{
			encryptionKey: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			value:      "\x00",
			variableID: -1,
		},
		{
			encryptionKey: []byte{
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			},
			value:      "\xFF",
			variableID: -1,
		},
		{
			encryptionKey: []byte("01234567890123456789012345678901"),
			value:         "hello, world",
			variableID:    math.MaxInt32,
		},
		{
			encryptionKey: []byte("01234567890123456789012345678901"),
			value:         "hello, world",
			variableID:    math.MinInt32,
		},
	}

	for _, testcase := range testcases {
		f.Add(testcase.encryptionKey, testcase.value, testcase.variableID)
	}

	f.Fuzz(func(t *testing.T, encryptionKey []byte, value string, variableID int) {
		ciphertext, err := Encrypt(encryptionKey, value, variableID)
		if err != nil {
			t.Skip()
		}
		_, err = Decrypt(encryptionKey, ciphertext, variableID)
		if err != nil {
			t.Skip()
		}
	})
}
