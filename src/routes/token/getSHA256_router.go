package tokens

import (
	"magic-server-2026/src/gen"
	"magic-server-2026/src/utils"
	"time"

	"github.com/gofiber/fiber/v3"
)

func GetSHA256Token(c fiber.Ctx) error {
	existing := c.Cookies("_sha256")

	if existing != "" {
		return c.JSON(fiber.Map{
			"_sha256": existing,
			"reuse":   true,
		})
	}

	raw, err := gen.GenerateSHA256AESGCM()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "fail"})
	}

	utils.RevokeCookie(c, "_sha256")
	utils.SetSecureCookie(c, "_sha256", raw, 5*time.Second)

	return c.JSON(fiber.Map{
		"_sha256": raw,
		"new":     true,
	})
}
