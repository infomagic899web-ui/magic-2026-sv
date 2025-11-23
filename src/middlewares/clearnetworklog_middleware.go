package middlewares

import "github.com/gofiber/fiber/v3"

func SilentErrorsMiddleware(c fiber.Ctx) error {
	err := c.Next()

	if err != nil {
		// fiber normally logs errors here
		// we intercept them and return a silent clean response
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		c.Set("Cache-Control", "no-store, max-age=15")
		c.Set("Pragma", "no-cache")
		c.Set("Expires", "0")
		// only return status, no message => no console spam
		return c.Status(code).SendString("")
	}

	return nil
}
