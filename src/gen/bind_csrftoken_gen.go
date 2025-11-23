package gen

import (
	"magic-server-2026/src/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateBindedCSRF(sessionID string) (string, error) {
	claims := models.BindCSRFClaims{
		SessionID: sessionID,
		Type:      "_bind_csrf",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(CSRFSecret)
}
