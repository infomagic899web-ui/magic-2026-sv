package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

// RoleFilterMiddleware checks if user role matches allowed roles
func RoleFilterMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Example: role extracted from JWT claims or context
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing user role",
			})
		}

		// Compare against allowed roles
		for _, role := range allowedRoles {
			if strings.EqualFold(userRole.(string), role) {
				return c.Next()
			}
		}

		// If no match
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}
}
