package gen

import (
	"magic-server-2026/src/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret for access tokens (use env variable in production)
var AccessSecret = []byte("super_secret_access_key")

// GenerateAccessToken creates a signed JWT access token
func GenerateAccessToken(userID, sessionID, role string) (string, error) {
	claims := &models.AccessTokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "magic899",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(AccessSecret)
}
