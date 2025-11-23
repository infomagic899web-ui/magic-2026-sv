package middlewares

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/microcosm-cc/bluemonday"
)

// SanitizeMiddleware sanitizes string fields in incoming JSON
func SanitizeMiddleware(c fiber.Ctx) error {
	if c.Method() != fiber.MethodPost && c.Method() != fiber.MethodPut && c.Method() != fiber.MethodPatch {
		return c.Next()
	}

	// Read raw body
	body := c.Body()
	if len(body) == 0 {
		return c.Next()
	}

	// Parse JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		// Not JSON, skip sanitization
		return c.Next()
	}

	// Create a strict policy
	policy := bluemonday.UGCPolicy() // User-Generated Content safe policy
	// Optionally customize allowed tags/attributes to match frontend whitelist
	policy.AllowAttrs("href", "title", "rel", "target").OnElements("a")
	policy.AllowAttrs("src", "alt", "width", "height").OnElements("img")

	// Recursively sanitize all string fields
	sanitizeMap(jsonData, policy)

	// Replace body with sanitized JSON
	safeBody, _ := json.Marshal(jsonData)
	c.Request().SetBody(safeBody)

	return c.Next()
}

// sanitizeMap recursively sanitizes string fields
func sanitizeMap(data map[string]interface{}, policy *bluemonday.Policy) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			data[key] = strings.TrimSpace(policy.Sanitize(v))
		case map[string]interface{}:
			sanitizeMap(v, policy)
		case []interface{}:
			sanitizeSlice(v, policy)
		}
	}
}

func sanitizeSlice(data []interface{}, policy *bluemonday.Policy) {
	for i, value := range data {
		switch v := value.(type) {
		case string:
			data[i] = strings.TrimSpace(policy.Sanitize(v))
		case map[string]interface{}:
			sanitizeMap(v, policy)
		case []interface{}:
			sanitizeSlice(v, policy)
		}
	}
}
