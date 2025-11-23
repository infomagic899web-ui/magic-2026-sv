package tokens

import (
	"magic-server-2026/src/gen"
	"magic-server-2026/src/utils"
	"time"

	"github.com/gofiber/fiber/v3"
)

func GetCSRFToken(c fiber.Ctx) error {
	existing := c.Cookies("_csrf")

	if existing != "" {
		return c.JSON(fiber.Map{
			"_csrf": existing,
			"reuse": true,
		})
	}

	raw := gen.GenerateCSRFToken(c)

	utils.RevokeCookie(c, "_csrf")
	utils.SetSecureCookie(c, "_csrf", raw, 5*time.Second)

	return c.JSON(fiber.Map{
		"_csrf": raw,
		"new":   true,
	})
}
