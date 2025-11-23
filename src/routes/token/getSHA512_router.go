package tokens

import (
	"magic-server-2026/src/gen"
	"magic-server-2026/src/utils"
	"time"

	"github.com/gofiber/fiber/v3"
)

func GetSHA512Token(c fiber.Ctx) error {
	existing := c.Cookies("_sha512")

	if existing != "" {
		return c.JSON(fiber.Map{
			"_sha512": existing,
			"reuse":   true,
		})
	}

	raw, err := gen.GenerateSHA512AESGCM()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "fail"})
	}

	utils.RevokeCookie(c, "_sha512")
	utils.SetSecureCookie(c, "_sha512", raw, 5*time.Second)

	return c.JSON(fiber.Map{
		"_sha512": raw,
		"new":     true,
	})
}
