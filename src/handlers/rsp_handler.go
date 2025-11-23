package handlers

import (
	"magic-server-2026/src/gen"
	"magic-server-2026/src/utils"

	"github.com/gofiber/fiber/v3"
)

func RotateRSToken(c fiber.Ctx) error {
	utils.RevokeCookie(c, "_rsp")
	rawToken := gen.GenerateRSToken(c)
	c.Set("X-RSP-Token", rawToken)
	return c.SendString("RSP Rotated")
}
