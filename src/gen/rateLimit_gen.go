package gen

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRateLimitSecret(size int) (string, error) {
	if size != 32 && size != 64 && size != 128 && size != 256 {
		return "", fmt.Errorf("invalid size: %d; allowed sizes are 32, 64, 128, 256", size)
	}

	secret := make([]byte, size)
	_, err := rand.Read(secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate secret: %v", err)
	}

	// Encode as URL-safe base64 string
	return base64.RawURLEncoding.EncodeToString(secret), nil
}
