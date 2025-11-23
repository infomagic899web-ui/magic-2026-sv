package gen

import (
	"magic899-server/src/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	CSRFSecret = []byte("super_secret_csrf_key")
	RSPSecret  = []byte("super_secret_rsp_key")
)

func GenerateBindedRSP(sessionID string) (string, error) {
	claims := models.BindRSPClaims{
		SessionID: sessionID,
		Type:      "_bind_rsp",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(RSPSecret)
}
