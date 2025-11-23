package models

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	UserID    string `json:"uid"`
	SessionID string `json:"sid"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID    string `json:"uid"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

type BindCSRFClaims struct {
	SessionID string `json:"sid"`
	Type      string `json:"type"` // "csrf"
	jwt.RegisteredClaims
}

type BindRSPClaims struct {
	SessionID string `json:"sid"`
	Type      string `json:"type"` // "rsp"
	jwt.RegisteredClaims
}
