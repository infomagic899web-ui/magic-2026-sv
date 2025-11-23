package middlewares

import (
	"magic-server-2026/src/gen"

	"github.com/gofiber/fiber/v3"
)

// BindTokenMiddleware validates bind_csrf and bind_rsp, checks revocation, and rotates tokens
func BindTokenMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing session_id")
		}

		csrfToken := c.Cookies("bind_csrf")
		rspToken := c.Cookies("bind_rsp")

		// Validate tokens with session binding
		csrfClaims, errCSRF := ValidateBindedCSRF(csrfToken, sessionID)
		rspClaims, errRSP := ValidateBindedRSP(rspToken, sessionID)

		// If either token is invalid or session mismatch, reject
		if errCSRF != nil || errRSP != nil ||
			csrfClaims.SessionID != sessionID ||
			rspClaims.SessionID != sessionID {
			return fiber.NewError(fiber.StatusForbidden, "invalid or revoked bind tokens")
		}

		// Rotate tokens immediately (5 seconds lifetime)
		newCSRF, _ := gen.GenerateBindedCSRF(sessionID)
		newRSP, _ := gen.GenerateBindedRSP(sessionID)

		c.Cookie(&fiber.Cookie{
			Name:     "bind_csrf",
			Value:    newCSRF,
			Path:     "/csrf",
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
			MaxAge:   5,
		})

		c.Cookie(&fiber.Cookie{
			Name:     "bind_rsp",
			Value:    newRSP,
			Path:     "/rsp",
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
			MaxAge:   5,
		})

		return c.Next()
	}
}
