package crypto

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

func Encrypt(encryptionKey []byte, value string, variableID int) ([]byte, error) {
	cipher, err := chacha20poly1305.NewX(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	nonce := make([]byte, cipher.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	var ad [4]byte
	binary.LittleEndian.PutUint32(ad[:], uint32(variableID))

	plaintext := []byte(value)
	ciphertext := make([]byte, 0, len(plaintext)+cipher.Overhead()+cipher.NonceSize())
	ciphertext = append(ciphertext, nonce[:]...)
	ciphertext = cipher.Seal(ciphertext, nonce, plaintext, ad[:])
	return ciphertext, nil
}

func Decrypt(encryptionKey []byte, ciphertext []byte, variableID int) ([]byte, error) {
	cipher, err := chacha20poly1305.NewX(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	var ad [4]byte
	binary.LittleEndian.PutUint32(ad[:], uint32(variableID))

	plaintextSize := len(ciphertext) - cipher.NonceSize() - cipher.Overhead()
	if plaintextSize < 0 {
		return nil, fmt.Errorf("ciphertext is too short")
	}
	plaintext := make([]byte, 0, plaintextSize)
	plaintext, err = cipher.Open(plaintext, ciphertext[:cipher.NonceSize()], ciphertext[cipher.NonceSize():], ad[:])
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}
	return plaintext, nil
}
