package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSessionToken() string {
	token := make([]byte, 32) // 32 bytes for token
	_, err := rand.Read(token)
	if err != nil {
		panic("Error generating session token")
	}
	return base64.URLEncoding.EncodeToString(token) // Encoding the token in base64 URL format
}
