package gen

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"magic-server-2026/src/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

var secretKey []byte

// GenerateCSRFToken creates an AES-GCM encrypted token and sets it as a cookie
func GenerateRSToken(c fiber.Ctx) string {
	// Create AES block cipher
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		fmt.Println("AES error:", err)
		return ""
	}

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("GCM error:", err)
		return ""
	}

	// Nonce (GCM standard = 12 bytes)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("Nonce error:", err)
		return ""
	}

	// Payload: session + timestamp
	payload := fmt.Sprintf("%s:%d", c.IP(), time.Now().UnixNano())

	// Encrypt
	ciphertext := aesGCM.Seal(nil, nonce, []byte(payload), nil)

	// Combine nonce + ciphertext
	token := append(nonce, ciphertext...)

	// Encode for transport
	encoded := base64.RawURLEncoding.EncodeToString(token)

	utils.RevokeCookie(c, "_rsp")

	// ✅ Set cookie (not HttpOnly so JS can read for AJAX requests, SameSite Strict recommended)
	utils.SetSecureCookie(c, "_rsp", encoded, 30*time.Minute)

	return encoded
}

func ValidateRSToken(c fiber.Ctx, token string) error {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return errors.New("invalid RSP encoding")
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return errors.New("AES init error")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return errors.New("GCM init error")
	}

	if len(raw) < aesGCM.NonceSize() {
		return errors.New("malformed RSP token")
	}

	nonce := raw[:aesGCM.NonceSize()]
	ciphertext := raw[aesGCM.NonceSize():]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return errors.New("RSP decryption failed")
	}

	// Parse payload: IP:timestamp
	parts := strings.Split(string(plaintext), ":")
	if len(parts) != 2 {
		return errors.New("invalid RSP payload format")
	}

	ip := parts[0]
	ts, err := time.ParseDuration(fmt.Sprintf("%sns", parts[1]))
	if err != nil {
		return errors.New("invalid RSP timestamp")
	}

	// Enforce short lifetime (10 seconds)
	if time.Since(time.Unix(0, int64(ts))) > 5*time.Minute {
		return errors.New("RSP expired")
	}

	// IP-bound check
	if ip != c.IP() {
		return errors.New("RSP token not from same client")
	}

	return nil
}
