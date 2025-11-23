package gen

import (
	"time"

	"magic899-server/src/models"

	"github.com/golang-jwt/jwt/v5"
)

// Secret for refresh tokens (use env variable in production)
var RefreshSecret = []byte("super_secret_refresh_key")

// GenerateRefreshToken creates a signed JWT refresh token
func GenerateRefreshToken(userID, sessionID string) (string, error) {
	claims := &models.RefreshTokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "magic899",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(RefreshSecret)
}
