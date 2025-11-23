package gen

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"magic-server-2026/src/utils"
	"time"

	"github.com/gofiber/fiber/v3"
)

func GenerateNonceToken(c fiber.Ctx) string {
	block, err := aes.NewCipher(secretKey)
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

	// Payload: IP + timestamp
	payload := fmt.Sprintf("%s:%d", c.IP(), time.Now().UnixNano())
	ciphertext := aesGCM.Seal(nil, nonce, []byte(payload), nil)

	// Combine nonce + ciphertext
	token := append(nonce, ciphertext...)
	encoded := base64.RawURLEncoding.EncodeToString(token)

	utils.RevokeCookie(c, "_nonce")
	utils.SetSecureCookie(c, "_nonce", encoded, 30*time.Minute) // short-lived 5s

	return encoded
}
