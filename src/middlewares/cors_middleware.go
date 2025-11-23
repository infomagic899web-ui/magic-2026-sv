package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"magic-server-2026/src/gen"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

// ⚠️ Secure secret from .env
var hmacSecret = []byte("replace-with-a-secure-random-hmac-secret-32bytes-min")

// signNonce creates an HMAC-signed token based on a nonce and session ID
func signNonce(nonce, sessionID string) string {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	mac := hmac.New(sha256.New, hmacSecret)
	mac.Write([]byte(nonce + "|" + sessionID + "|" + ts))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s.%s.%s", nonce, ts, sig)
}

// CORSMiddleware aligns backend CSP with frontend, including SHA256 & SHA512
func CORSMiddleware(c fiber.Ctx) error {
	env := os.Getenv("ENV")
	isProd := env == "production"

	// Determine allowed frontend origin
	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		if isProd {
			frontendOrigin = "https://demo.magic899.com"
		} else {
			frontendOrigin = "http://localhost:5179"
		}
	}

	origin := c.Get("Origin")
	allowOrigin := frontendOrigin
	if !isProd && origin != "" && (strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1")) {
		allowOrigin = origin
	}

	// Generate nonce + signed token
	sessionID := c.Cookies("session_id")
	nonce := gen.GenerateNonceToken(c)
	token := signNonce(nonce, sessionID)

	// Example: dynamic SHA256 & SHA512 hashes for inline scripts/styles
	sha256sum := c.Cookies("_sha256")
	sha512sum := c.Cookies("_sha512")

	// CORS headers
	c.Set("Access-Control-Allow-Origin", allowOrigin)
	c.Set("Access-Control-Allow-Credentials", "true")
	c.Set("Access-Control-Allow-Headers", strings.Join([]string{
		"Content-Type",
		"Authorization",
		"X-CSRF-Token",
		"X-RSP-Token",
		"X-Nonce",
		"X-Nonce-Token",
		"X-Hash256",
		"X-Hash512",
		"Accept",
		"Origin",
		"Cache-Control",
	}, ", "))
	c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	c.Set("Cache-Control", "no-store, max-age=15")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	// Security headers
	c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	c.Set("Cross-Origin-Opener-Policy", "same-origin")
	c.Set("Cross-Origin-Embedder-Policy", "unsafe-none")
	c.Set("Cross-Origin-Resource-Policy", "cross-origin")
	c.Set("X-Frame-Options", "SAMEORIGIN")
	c.Set("X-XSS-Protection", "1; mode=block")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Set("Permissions-Policy", "geolocation=(), camera=(), microphone=(), fullscreen=(self)")

	// Nonce headers
	c.Set("X-Nonce", nonce)
	c.Set("X-Nonce-Token", token)

	backendOrigin := os.Getenv("SERVER_ORIGIN")

	// Content-Security-Policy including SHA256 & SHA512
	csp := fmt.Sprintf(strings.Join([]string{
		"default-src 'self'",
		fmt.Sprintf("script-src 'self' 'nonce-%s' 'strict-dynamic' 'sha256-%s' 'sha512-%s'", nonce, sha256sum, sha512sum),
		fmt.Sprintf("style-src 'self' 'nonce-%s' 'unsafe-hashes' 'sha256-%s' 'sha512-%s'", nonce, sha256sum, sha512sum),
		"img-src 'self' data: https://i.ytimg.com https://encrypted-tbn0.gstatic.com https://media.philstar.com https://i.scdn.co blob: https://cdn.weatherapi.com https://www.youtube.com https://www.youtube-nocookie.com",
		"font-src 'self' https://fonts.gstatic.com https://cdnjs.cloudflare.com",
		fmt.Sprintf("connect-src 'self' https://demo.magic899.com https://apidemo.magic899.com https://magic-89-9-2026-kzgl376xi-renstrio24ps-projects.vercel.app https://api.weatherapi.com https://cdn.weatherapi.com https://www.youtube.com %s", frontendOrigin),
		"frame-src https://www.youtube.com https://www.youtube-nocookie.com",
		"frame-ancestors 'none'",
		"object-src 'none'",
		"base-uri 'self'",
		"form-action 'self'",
		"manifest-src 'self'",
		"worker-src 'self'",
		"trusted-types enforcedPolicy enforcedReactPolicy",
		"upgrade-insecure-requests",
		fmt.Sprintf("report-uri %s/api/enforce-csp-report", backendOrigin),
	}, "; "))
	c.Set("Content-Security-Policy", csp)

	// Preflight
	if c.Method() == fiber.MethodOptions {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.Next()
}
