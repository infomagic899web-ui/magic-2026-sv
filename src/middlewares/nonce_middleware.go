package middlewares

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

func NonceMiddleware(c fiber.Ctx) error {
	// Generate cryptographically random 16-byte nonce
	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate nonce")
	}
	nonce := base64.RawURLEncoding.EncodeToString(nonceBytes)

	// Sign the nonce with timestamp using HMAC
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	mac := hmac.New(sha256.New, hmacSecret)
	mac.Write([]byte(nonce + "|" + ts))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	token := nonce + "." + ts + "." + sig

	// Set short-lived nonce cookie (readable by SSR)
	c.Cookie(&fiber.Cookie{
		Name:     "_nonce",
		Value:    nonce,
		Path:     "/",
		Secure:   true,
		HTTPOnly: false, // frontend can read it
		SameSite: "Lax",
		MaxAge:   10, // 10 seconds
	})

	// Expose to frontend & CSP
	c.Set("X-Nonce", nonce)
	c.Set("X-RSP-Nonce", token)
	c.Locals("nonce", nonce)

	return c.Next()
}

// --- Middleware: verify nonce token ---
func VerifyNonceTokenMiddleware(c fiber.Ctx) error {
	nonce := c.Get("X-Nonce")
	token := c.Get("X-RSP-Nonce")

	if nonce == "" || token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "missing nonce or token",
		})
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "invalid nonce token format",
		})
	}

	tokenNonce, tsStr, sig := parts[0], parts[1], parts[2]
	if tokenNonce != nonce {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "nonce mismatch",
		})
	}

	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "invalid timestamp",
		})
	}

	// Expire after 10 seconds
	if time.Now().Unix()-ts > 10 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "expired nonce token",
		})
	}

	// Validate HMAC (without session_id)
	mac := hmac.New(sha256.New, hmacSecret)
	mac.Write([]byte(nonce + "|" + tsStr))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedSig), []byte(sig)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "invalid nonce signature",
		})
	}

	// ✅ Valid token
	return c.Next()
}
