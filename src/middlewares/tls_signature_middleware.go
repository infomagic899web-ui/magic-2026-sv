package middlewares

import "github.com/gofiber/fiber/v3"

func TLSRedirectMiddleware(c fiber.Ctx) error {
	if c.Protocol() != "https" {
		c.Redirect().To("https://" + c.Hostname() + c.OriginalURL())
		return nil
	}
	return c.Next()
}
