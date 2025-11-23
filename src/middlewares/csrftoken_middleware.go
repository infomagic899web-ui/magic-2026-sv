package middlewares

import (
	"magic-server-2026/src/gen"
	"os"

	"github.com/gofiber/fiber/v3"
)

/*
	CSRF Token Middleware
	- Prevents Data Modification
	- Checks if the CSRF token is valid (encrypted + time-limited)
	- Protects POST, PUT, PATCH, DELETE requests
*/

var isProd = os.Getenv("ENV") == "production"

func CSRFTokenMiddleware(c fiber.Ctx) error {
	// Only protect state-changing requests
	if c.Method() != fiber.MethodPost &&
		c.Method() != fiber.MethodPut &&
		c.Method() != fiber.MethodPatch &&
		c.Method() != fiber.MethodDelete {
		return c.Next()
	}

	csrfCookie := c.Cookies("_csrf")
	csrfHeader := c.Get("X-CSRF-Token")

	if csrfCookie == "" || csrfHeader == "" {
		if isProd {
			return fiber.ErrForbidden
		}
		return fiber.NewError(fiber.StatusForbidden, "Missing CSRF token - Retrying Again...")
	}
	// Basic equality check first
	if csrfCookie != csrfHeader {
		if isProd {
			return fiber.ErrForbidden
		}
		return fiber.NewError(fiber.StatusForbidden, "CSRF token mismatch - Retrying Again...")
	}

	// Deep validation (AES-GCM decrypt + expiry + IP binding)
	if err := gen.ValidateCSRFToken(c, csrfCookie); err != nil {
		if isProd {
			return fiber.ErrForbidden
		}
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	return c.Next()
}
