package gen

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"io"
)

// GenerateSHA512AESGCM generates a SHA512 hash, encrypts it with AES-GCM, and returns base64 string.
func GenerateSHA512AESGCM() (string, error) {
	// Step 1: generate random 64-byte input
	randomBytes := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		return "", err
	}

	// Step 2: hash with SHA512
	hashed := sha512.Sum512(randomBytes)

	// Step 3: AES-GCM encrypt the hash
	key := make([]byte, 32) // 256-bit AES key
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, hashed[:], nil)

	// Step 4: Return base64 string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
