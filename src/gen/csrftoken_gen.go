package gen

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"magic-server-2026/src/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

var (
	mu sync.RWMutex
)

// Init starts secret regeneration every 15 seconds
func Init() {
	// Initial secret
	generateSecret()

	// Regenerate secret every 15 seconds
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			generateSecret()
		}
	}()
}

func generateSecret() {
	newSecret := make([]byte, 32)
	if _, err := rand.Read(newSecret); err != nil {
		log.Fatal("Failed to generate new secret:", err)
	}
	mu.Lock()
	secretKey = newSecret
	mu.Unlock()
	log.Println("✅ Military Secret regenerated")
}

// getSecret safely returns the current secret
func getSecret() []byte {
	mu.RLock()
	defer mu.RUnlock()
	secretCopy := make([]byte, len(secretKey))
	copy(secretCopy, secretKey)
	return secretCopy
}

func GenerateCSRFToken(c fiber.Ctx) string {
	secret := getSecret()
	block, err := aes.NewCipher(secret)
	if err != nil {
		log.Println("AES error:", err)
		return ""
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("GCM error:", err)
		return ""
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("Nonce error:", err)
		return ""
	}

	// Payload: IP + timestamp (Unix Nano)
	payload := fmt.Sprintf("%s:%d", c.IP(), time.Now().UnixNano())
	ciphertext := aesGCM.Seal(nil, nonce, []byte(payload), nil)

	token := append(nonce, ciphertext...)
	encoded := base64.RawURLEncoding.EncodeToString(token)

	// Revoke previous CSRF and issue a new one for 30 minutes
	utils.RevokeCookie(c, "_csrf")
	utils.SetSecureCookie(c, "_csrf", encoded, 30*time.Minute)

	return encoded
}

// ValidateCSRFToken checks token validity, IP, and expiration
func ValidateCSRFToken(c fiber.Ctx, token string) error {
	secret := getSecret()
	block, err := aes.NewCipher(secret)
	if err != nil {
		return errors.New("AES init error")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return errors.New("GCM init error")
	}

	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return errors.New("invalid token encoding")
	}

	nonceSize := aesGCM.NonceSize()
	if len(decoded) < nonceSize {
		return errors.New("token too short")
	}

	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("decryption failed")
	}

	parts := strings.Split(string(plaintext), ":")
	if len(parts) != 2 {
		return errors.New("invalid payload format")
	}

	tokenIP := parts[0]
	timestampStr := parts[1]

	ts, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return errors.New("invalid timestamp")
	}

	// Updated expiration: 30 minutes
	if time.Since(time.Unix(0, ts)) > 30*time.Minute {
		return errors.New("CSRF token expired")
	}

	// Optional: IP binding check
	if tokenIP != c.IP() {
		return errors.New("IP mismatch")
	}

	return nil
}
