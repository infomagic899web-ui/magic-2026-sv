package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateNonce() string {
	// Generate a 16-byte random value
	nonce := make([]byte, 16)
	_, err := rand.Read(nonce)
	if err != nil {
		panic("Failed to generate nonce")
	}
	return base64.StdEncoding.EncodeToString(nonce)
}
