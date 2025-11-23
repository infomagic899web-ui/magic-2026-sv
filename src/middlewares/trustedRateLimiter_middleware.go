package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

// Configurable values
var (
	trustedDomains = []string{
		"https://magic899.com",
		"https://app.magic899.com",
		"https://demo-test.magic899.com",
		"https://cloudflare.com",
	}

	untrustedWindow      = 5 * time.Hour
	rateLimiterSecretEnv = "RATE_LIMIT_SECRET"
)

// entry holds per-fingerprint rate-limiter state
type entry struct {
	count  int
	expiry time.Time
}

// RateLimiter is the in-memory store + mutex
type RateLimiter struct {
	mtx             sync.Mutex
	store           map[string]*entry
	cleanupInterval time.Duration
}

// NewRateLimiter creates a new limiter
func NewRateLimiter(cleanupInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		store:           make(map[string]*entry),
		cleanupInterval: cleanupInterval,
	}
	go rl.cleanupLoop()
	return rl
}

func (r *RateLimiter) cleanupLoop() {
	t := time.NewTicker(r.cleanupInterval)
	for range t.C {
		now := time.Now()
		r.mtx.Lock()
		for k, v := range r.store {
			if v.expiry.Before(now) {
				delete(r.store, k)
			}
		}
		r.mtx.Unlock()
	}
}

func (r *RateLimiter) allowOnce(key string, window time.Duration) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	now := time.Now()
	e, ok := r.store[key]
	if !ok || e.expiry.Before(now) {
		// Allow first request
		r.store[key] = &entry{
			count:  1,
			expiry: now.Add(window),
		}
		return true
	}

	if e.count < 1 && e.expiry.After(now) {
		e.count++
		return true
	}

	return false
}

// Check if origin is trusted
func isTrustedOrigin(originHeader string) bool {
	originHeader = strings.TrimSpace(originHeader)
	if originHeader == "" {
		return false
	}

	// Allow localhost for dev
	if strings.Contains(originHeader, "localhost") || strings.Contains(originHeader, "127.0.0.1") {
		return true
	}

	// Check explicit trusted domains
	for _, td := range trustedDomains {
		if strings.EqualFold(td, originHeader) || strings.HasPrefix(originHeader, td) {
			return true
		}
	}

	// Trust any render.com subdomain
	if strings.HasSuffix(originHeader, ".onrender.com") {
		return true
	}

	return false
}

// Build fingerprint string from request
func buildFingerprint(c fiber.Ctx) string {
	ip := c.IP()
	ua := c.Get("User-Agent", "")
	acceptLang := c.Get("Accept-Language", "")
	secChUa := c.Get("Sec-CH-UA", "")
	origin := c.Get("Origin", "")
	referer := c.Get("Referer", "")
	deviceID := c.Get("X-Device-ID", "")   // optional for mobile
	appBundle := c.Get("X-App-Bundle", "") // optional iOS WebView

	parts := []string{
		ip,
		ua,
		acceptLang,
		secChUa,
		origin,
		referer,
		deviceID,
		appBundle,
		c.Path(),
		c.Method(),
	}

	return strings.Join(parts, "|")
}

// HMAC fingerprint with secret from env
func hmacFingerprint(secret, fingerprint string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(fingerprint))
	if err != nil {
		return "", err
	}
	sum := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(sum), nil
}

// Middleware factory
func MakeTrustedRateLimiterMiddleware(rl *RateLimiter) fiber.Handler {
	secret := os.Getenv(rateLimiterSecretEnv)
	if secret == "" {
		log.Printf("[WARN] %s is empty - HMAC fingerprint will be weaker (set env var)", rateLimiterSecretEnv)
	}

	return func(c fiber.Ctx) error {
		origin := c.Get("Origin", "")
		referer := c.Get("Referer", "")
		appBundle := c.Get("X-App-Bundle", "")

		// Trusted domains bypass
		if isTrustedOrigin(origin) || isTrustedOrigin(referer) {
			return c.Next()
		}

		// Optional iOS app bundle bypass
		if appBundle == "com.magic899.app" {
			return c.Next()
		}

		// Untrusted request: fingerprint + HMAC
		fp := buildFingerprint(c)
		fpHmac := fp
		if secret != "" {
			if h, err := hmacFingerprint(secret, fp); err == nil {
				fpHmac = h
			}
		} else {
			fpHmac = base64.RawURLEncoding.EncodeToString([]byte(fp))
		}

		allowed := rl.allowOnce(fpHmac, untrustedWindow)
		if !allowed {
			// Log full details internally
			log.Printf(
				"[RATE-THROTTLE] blocked untrusted request - origin=%s referer=%s ip=%s path=%s ua=%s fingerprint=%s",
				origin, referer, c.IP(), c.Path(), c.Get("User-Agent", ""), fpHmac,
			)

			// Compute Retry-After
			r := 60
			rl.mtx.Lock()
			if e, ok := rl.store[fpHmac]; ok {
				if ttl := int(time.Until(e.expiry).Seconds()); ttl > 0 {
					r = ttl
				}
			}
			rl.mtx.Unlock()

			// Return generic hidden message
			c.Set("Retry-After", strconv.Itoa(r))
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "blocked",
				"message":     "request blocked", // hidden message
				"retry_after": r,
			})
		}

		// Log allowed untrusted request internally
		log.Printf(
			"[RATE-ALLOW] untrusted-first-request origin=%s referer=%s ip=%s path=%s ua=%s fingerprint=%s",
			origin, referer, c.IP(), c.Path(), c.Get("User-Agent", ""), fpHmac,
		)

		return c.Next()
	}
}
