package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"
)

var signedURLSecret = []byte("YOUR_SECRET_KEY") // Change to a strong secret in env

// GenerateSignedURL generates a temporary signed URL for a filename
func GenerateSignedURL(filename string, duration time.Duration) (string, error) {
	expires := time.Now().Add(duration).Unix()
	data := fmt.Sprintf("%s:%d", filename, expires)

	h := hmac.New(sha256.New, signedURLSecret)
	h.Write([]byte(data))
	sig := hex.EncodeToString(h.Sum(nil))

	// URL-encode filename
	escapedFilename := url.PathEscape(filename)

	// Return query string URL: /player/:filename?expires=...&sig=...
	return fmt.Sprintf("/api/v1/player/%s?expires=%d&sig=%s", escapedFilename, expires, sig), nil
}

// ValidateSignedURL validates the HMAC signature and expiration
func ValidateSignedURL(filename string, expires int64, sig string) bool {
	data := fmt.Sprintf("%s:%d", filename, expires)

	h := hmac.New(sha256.New, signedURLSecret)
	h.Write([]byte(data))
	expected := hex.EncodeToString(h.Sum(nil))

	// Check signature match and expiration
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return false
	}

	if time.Now().Unix() > expires {
		return false
	}

	return true
}
