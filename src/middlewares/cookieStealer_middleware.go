package middlewares

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v3"
)

var (
	// Precompile and anchor regex (avoids unnecessary substring scanning)
	cookieRegex = regexp.MustCompile(`(?i)\bdocument\.cookie\b`)
	xssKeywords = []*regexp.Regexp{
		regexp.MustCompile(`(?i)\bwindow\.location\b`),
		regexp.MustCompile(`(?i)\bfetch\(`),
		regexp.MustCompile(`(?i)\bXMLHttpRequest\b`),
	}
)

// CookieStealBlocker blocks suspicious requests containing cookie-stealing patterns.
func CookieStealBlocker() fiber.Handler {
	return func(c fiber.Ctx) error {
		uri := string(c.Request().RequestURI())

		// ✅ 1. Fast path: if no suspicious substring, skip regex
		if !containsSuspiciousSubstrings(uri) {
			return c.Next()
		}

		// ✅ 2. Deep inspection
		if cookieRegex.MatchString(uri) {
			return forbidden(c, "document.cookie in URI")
		}
		for _, re := range xssKeywords {
			if re.MatchString(uri) {
				return forbidden(c, "potential XSS keyword in URI")
			}
		}

		// ✅ 3. Check query params (defense-in-depth)
		query := c.Request().URI().QueryArgs().String()
		if cookieRegex.MatchString(query) {
			return forbidden(c, "document.cookie in query")
		}

		return c.Next()
	}
}

func containsSuspiciousSubstrings(s string) bool {
	return (len(s) > 0 &&
		(strings.Contains(s, "document") ||
			strings.Contains(s, "cookie") ||
			strings.Contains(s, "window") ||
			strings.Contains(s, "fetch") ||
			strings.Contains(s, "XMLHttpRequest")))
}

func forbidden(c fiber.Ctx, reason string) error {
	// Don’t log sensitive payloads in production
	if c.App().Config().AppName != "prod" {
		println("🚨 CookieStealBlocker:", reason)
	}
	return c.SendStatus(fiber.StatusForbidden)
}
