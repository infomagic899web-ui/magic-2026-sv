package tokens

import (
	"magic-server-2026/src/gen"
	"magic-server-2026/src/utils"
	"time"

	"github.com/gofiber/fiber/v3"
)

func GetNonceToken(c fiber.Ctx) error {
	existing := c.Cookies("_nonce")

	if existing != "" {
		return c.JSON(fiber.Map{
			"_nonce": existing,
			"reuse":  true,
		})
	}

	raw := gen.GenerateNonceToken(c)

	utils.RevokeCookie(c, "_nonce")
	utils.SetSecureCookie(c, "_nonce", raw, 5*time.Second)

	return c.JSON(fiber.Map{
		"_nonce": raw,
		"new":    true,
	})
}
