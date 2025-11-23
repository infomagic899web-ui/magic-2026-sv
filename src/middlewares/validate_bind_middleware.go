package middlewares

import (
	"errors"
	"magic-server-2026/src/models"

	"github.com/golang-jwt/jwt/v5"
)

// Secrets (replace with env variables in production)
var (
	CSRFSecret    = []byte("super_secret_csrf_key")
	RSPSecret     = []byte("super_secret_rsp_key")
	RefreshSecret = []byte("super_secret_refresh_key")
	AccessSecret  = []byte("super_secret_access_key")
)

// --------------------
// BIND CSRF
// --------------------
func ValidateBindedCSRF(tokenStr string, sessionID string) (*models.BindCSRFClaims, error) {
	if TokenRevoker.IsCSRFRevoked(tokenStr) {
		return nil, errors.New("CSRF token revoked")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &models.BindCSRFClaims{}, func(t *jwt.Token) (interface{}, error) {
		return CSRFSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.BindCSRFClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid CSRF token")
	}

	if claims.SessionID != sessionID {
		return nil, errors.New("CSRF token session mismatch")
	}

	// Revoke old token after validation
	TokenRevoker.RevokeCSRF(tokenStr)

	return claims, nil
}

// --------------------
// BIND RSP
// --------------------
func ValidateBindedRSP(tokenStr string, sessionID string) (*models.BindRSPClaims, error) {
	if TokenRevoker.IsRSPRevoked(tokenStr) {
		return nil, errors.New("RSP token revoked")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &models.BindRSPClaims{}, func(t *jwt.Token) (interface{}, error) {
		return RSPSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.BindRSPClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid RSP token")
	}

	if claims.SessionID != sessionID {
		return nil, errors.New("RSP token session mismatch")
	}

	// Revoke old token after validation
	TokenRevoker.RevokeRSP(tokenStr)

	return claims, nil
}

// --------------------
// REFRESH TOKEN
// --------------------
func ValidateRefreshToken(tokenStr string) (*models.RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return RefreshSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}

// --------------------
// ACCESS TOKEN
// --------------------
func ValidateAccessToken(tokenStr string) (*models.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return AccessSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}
