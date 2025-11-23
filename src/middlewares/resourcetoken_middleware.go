package middlewares

import (
	"magic-server-2026/src/gen"

	"github.com/gofiber/fiber/v3"
)

func ResourceTokenMiddleware(c fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return c.Next()
	}

	// allow token generation endpoints
	path := c.Path()
	if path == "/api/token/rsp" ||
		path == "/api/token/csrf" ||
		path == "/api/token/nonce" ||
		path == "/api/token/sha256" ||
		path == "/api/token/sha512" {
		return c.Next()
	}

	rsToken := c.Cookies("_rsp")
	clientRSToken := c.Get("X-RSP-Token")

	if rsToken == "" || clientRSToken == "" {
		if isProd {
			return c.Status(fiber.StatusForbidden).SendString("")
		}
		return c.Status(fiber.StatusForbidden).SendString("Missing RSP token")
	}

	if rsToken != clientRSToken {
		if isProd {
			return c.Status(fiber.StatusForbidden).SendString("")
		}
		return c.Status(fiber.StatusForbidden).SendString("Forbidden")
	}

	if err := gen.ValidateRSToken(c, rsToken); err != nil {
		if isProd {
			return c.Status(fiber.StatusForbidden).SendString("")
		}
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	return c.Next()
}
