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

