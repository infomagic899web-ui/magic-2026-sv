package gen

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// Secret key for AES-GCM (32 bytes for AES-256)
// In production, use an env variable
var AESKey = []byte("super_secret_32_byte_key_1234567890") // must be 32 bytes

// GenerateSecureSessionID generates a random session ID and encrypts it with AES-GCM
func GenerateSecureSessionID() (string, error) {
	// Generate random 32-byte session ID
	plain := make([]byte, 32)
	if _, err := rand.Read(plain); err != nil {
		return "", err
	}

	// Create AES cipher
	block, err := aes.NewCipher(AESKey)
	if err != nil {
		return "", err
	}

	// Wrap in GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the session ID
	ciphertext := aesGCM.Seal(nonce, nonce, plain, nil)

	// Return hex string
	return hex.EncodeToString(ciphertext), nil
}

// DecryptSecureSessionID decrypts a session ID encrypted with AES-GCM
func DecryptSecureSessionID(encHex string) ([]byte, error) {
	data, err := hex.DecodeString(encHex)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(AESKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, err
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plain, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plain, nil
}
