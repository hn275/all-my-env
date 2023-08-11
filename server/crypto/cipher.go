package crypto

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/hn275/envhub/server/lib"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	VariableKey  []byte
	UserTokenKey []byte
	UserIDKey    []byte
)

func init() {
	VariableKey = getKey("VARIABLE_KEY")
	UserTokenKey = getKey("USER_TOKEN_KEY")
	UserIDKey = getKey("USER_ID_KEY")
}

// one issue
// need variableID to encrypt, this can be the repo id or something
func Encrypt(key, plaintext, ad []byte) ([]byte, error) {
	cipher, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	nonce := make([]byte, cipher.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := make([]byte, 0, len(plaintext)+cipher.Overhead()+cipher.NonceSize())
	ciphertext = append(ciphertext, nonce[:]...)
	ciphertext = cipher.Seal(ciphertext, nonce, plaintext, ad)
	return ciphertext, nil
}

func Decrypt(key, ciphertext, ad []byte) ([]byte, error) {
	cipher, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	plaintextSize := len(ciphertext) - cipher.NonceSize() - cipher.Overhead()
	if plaintextSize < 0 {
		return nil, fmt.Errorf("ciphertext is too short")
	}
	plaintext := make([]byte, 0, plaintextSize)
	plaintext, err = cipher.Open(
		plaintext,
		ciphertext[:cipher.NonceSize()],
		ciphertext[cipher.NonceSize():],
		ad,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}
	return plaintext, nil
}

func getKey(k string) []byte {
	v := []byte(lib.Getenv(k))
	if len(v) != 32 {
		fmt.Fprintf(
			os.Stderr,
			"invalid CIPHER_KEY length.\nExpected 32, got: %d",
			len(VariableKey),
		)
		os.Exit(1)
	}
	return v
}
