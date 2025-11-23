package utils

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

func SetSecureCookie(c fiber.Ctx, name, value string, expires time.Duration) {
	c.Cookie(&fiber.Cookie{
		Name:        name,
		Value:       value,
		Path:        "/",
		Expires:     time.Now().Add(expires), // <-- expire in future
		MaxAge:      int(expires.Seconds()),
		Secure:      true,
		HTTPOnly:    true,
		SameSite:    "Lax",
		Partitioned: true,
	})
}

func RevokeCookie(c fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name:        name,
		Value:       "",
		Path:        "/",
		Expires:     time.Now().Add(-1 * time.Hour),
		MaxAge:      0,
		Secure:      true,
		HTTPOnly:    true,
		SameSite:    "Lax",
		Partitioned: false,
	})
}
