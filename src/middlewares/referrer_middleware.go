package middlewares

import (
	"net/url"

	"github.com/gofiber/fiber/v3"
)

func ReferrerMiddleware(c fiber.Ctx) error {
	origin := c.Get("Origin")
	referer := c.Get("Referer")

	// PROD strict domain lock
	allowedProd := "https://demo.magic899.com"

	// DEV allow localhost (vite / react)
	allowedDev := "http://localhost:5179"

	// resolve allowedOrigin by env
	allowed := allowedProd
	if !isProd {
		allowed = allowedDev
	}

	// allow OPTIONS preflight
	if c.Method() == fiber.MethodOptions {
		return c.Next()
	}

	// block requests without Origin + Referer (typing directly in browser / curl)
	if origin == "" && referer == "" {
		return fiber.ErrForbidden
	}

	// validate Origin
	if origin != "" && origin != allowed {
		return fiber.ErrForbidden
	}

	// validate Referer
	if referer != "" {
		u, err := url.Parse(referer)
		if err != nil || (u.Scheme+"://"+u.Host) != allowed {
			return fiber.ErrForbidden
		}
	}

	return c.Next()
}
