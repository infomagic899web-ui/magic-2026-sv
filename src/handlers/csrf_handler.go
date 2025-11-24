package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

var csrfStore = struct {
	sync.RWMutex
	tokens map[string]time.Time
}{
	tokens: make(map[string]time.Time),
}

// Generate random CSRF token
func generateCSRFToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Issue new CSRF token (GET /csrf-token)
func CSRFTokenHandler(c fiber.Ctx) error {
	token, err := generateCSRFToken()
	if err != nil {
		log.Println("[CSRF] Error generating token:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// store token for 30 minutes
	csrfStore.Lock()
	csrfStore.tokens[token] = time.Now().Add(30 * time.Minute)
	csrfStore.Unlock()

	// return token in header
	c.Set("X-CSRF-Token", token)

	return c.SendStatus(fiber.StatusNoContent)
}
