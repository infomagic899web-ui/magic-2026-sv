package tokens

import (
	"magic-server-2026/src/gen"
	"magic-server-2026/src/utils"
	"time"

	"github.com/gofiber/fiber/v3"
)

func GetRSToken(c fiber.Ctx) error {
	existing := c.Cookies("_rsp")

	if existing != "" {
		return c.JSON(fiber.Map{
			"_rsp":  existing,
			"reuse": true,
		})
	}

	raw := gen.GenerateRSToken(c)

	utils.RevokeCookie(c, "_rsp")
	utils.SetSecureCookie(c, "_rsp", raw, 5*time.Second)

	return c.JSON(fiber.Map{
		"_rsp": raw,
		"new":  true,
	})
}
