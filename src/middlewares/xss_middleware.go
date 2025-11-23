package middlewares

import (
	"regexp"

	"github.com/gofiber/fiber/v3"
)

/*
	XSS Middleware

	Purpose:
	- Blocks requests containing common XSS payload patterns.
	- Detects suspicious input in URI/query strings before it reaches handlers.
	- Provides an additional defense layer against reflected/stored XSS.

	Note:
	- Use with proper output encoding, CSP, and input validation.
	- Regex filters are heuristic and may cause false positives.
*/

// Precompile common XSS regex patterns (case-insensitive)
var xssPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)<script.*?>.*?</script.*?>`), // <script> tags
	regexp.MustCompile(`(?i)on\w+=`),                     // inline event handlers: onerror=, onclick=, etc.
	regexp.MustCompile(`(?i)javascript:`),                // javascript: URIs
	regexp.MustCompile(`(?i)eval\(`),                     // eval(
	regexp.MustCompile(`(?i)alert\(`),                    // alert(
	regexp.MustCompile(`(?i)<.*?iframe.*?>`),             // <iframe>
	regexp.MustCompile(`(?i)<.*?img.*?src=.*?>`),         // <img src=...>
}

// XSSBlocker blocks requests with suspicious XSS patterns
func XSSBlocker() fiber.Handler {
	return func(c fiber.Ctx) error {
		uri := string(c.Request().RequestURI())

		for _, pattern := range xssPatterns {
			if pattern.MatchString(uri) {
				return c.SendStatus(fiber.StatusForbidden) // 403
			}
		}

		return c.Next()
	}
}
