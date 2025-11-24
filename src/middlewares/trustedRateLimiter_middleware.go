package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"magic-server-2026/src/gen"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

var (
	trustedDomains = []string{
		"magic899.com",
		"demo-test.magic899.com",
	}

	untrustedWindow = 5 * time.Hour
	trustCloudflare = true
)

type entry struct {
	expiry time.Time
}

type RateLimiter struct {
	mtx             sync.Mutex
	store           map[string]*entry
	cleanupInterval time.Duration
}

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
	defer t.Stop()
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
		r.store[key] = &entry{
			expiry: now.Add(window),
		}
		return true
	}
	return false
}

func isTrustedHost(originHeader string) bool {
	if originHeader == "" {
		return false
	}

	u, err := url.Parse(originHeader)
	host := originHeader
	if err == nil {
		host = u.Hostname()
	}
	host = strings.ToLower(host)

	if host == "localhost" || host == "127.0.0.1" {
		return true
	}

	for _, td := range trustedDomains {
		if host == td || strings.HasSuffix(host, "."+td) {
			return true
		}
	}

	if strings.HasSuffix(host, ".onrender.com") {
		return true
	}

	return false
}

var cloudflareRanges []*net.IPNet

func init() {
	for _, cidr := range []string{
		"173.245.48.0/20",
		"103.21.244.0/22",
		"103.22.200.0/22",
		"103.31.4.0/22",
		"141.101.64.0/18",
		"108.162.192.0/18",
		"190.93.240.0/20",
		"188.114.96.0/20",
		"197.234.240.0/22",
		"198.41.128.0/17",
		"162.158.0.0/15",
		"104.16.0.0/13",
		"104.24.0.0/14",
		"172.64.0.0/13",
		"131.0.72.0/22",
	} {
		_, netw, err := net.ParseCIDR(cidr)
		if err == nil {
			cloudflareRanges = append(cloudflareRanges, netw)
		}
	}
}

func isCloudflareIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	for _, r := range cloudflareRanges {
		if r.Contains(ip) {
			return true
		}
	}
	return false
}

func preferRealIP(c fiber.Ctx) string {
	if trustCloudflare {
		cfIP := c.Get("CF-Connecting-IP")
		if cfIP != "" && isCloudflareIP(c.IP()) {
			return cfIP
		}
	}
	return c.IP()
}

func buildFingerprint(c fiber.Ctx) string {
	parts := []string{
		preferRealIP(c),
		c.Get("User-Agent", ""),
		c.Get("Accept-Language", ""),
		c.Get("Sec-CH-UA", ""),
		c.Get("Origin", ""),
		c.Get("Referer", ""),
		c.Get("X-Device-ID", ""),
		c.Get("X-App-Bundle", ""),
		c.Path(),
		c.Method(),
	}
	return strings.Join(parts, "|")
}

func hmacFingerprint(secret, fingerprint string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(fingerprint))
	sum := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(sum)
}

// Middleware factory
func MakeTrustedRateLimiterMiddleware(rl *RateLimiter) fiber.Handler {
	// Generate secret once in memory
	secret, err := gen.GenerateRateLimitSecret(128)
	if err != nil {
		log.Printf("[WARN] could not generate in-memory RATE_LIMIT_SECRET: %v", err)
		secret = ""
	}

	return func(c fiber.Ctx) error {
		origin := c.Get("Origin", "")
		referer := c.Get("Referer", "")
		appBundle := c.Get("X-App-Bundle", "")

		// Trusted bypass
		if isTrustedHost(origin) || isTrustedHost(referer) || appBundle == "com.magic899.app" {
			return c.Next()
		}

		fp := buildFingerprint(c)
		key := fp
		if secret != "" {
			key = hmacFingerprint(secret, fp)
		} else {
			h := sha256.Sum256([]byte(fp))
			key = base64.RawURLEncoding.EncodeToString(h[:])
		}

		allowed := rl.allowOnce(key, untrustedWindow)
		if !allowed {
			r := 60
			rl.mtx.Lock()
			if e, ok := rl.store[key]; ok {
				if ttl := int(time.Until(e.expiry).Seconds()); ttl > 0 {
					r = ttl
				}
			}
			rl.mtx.Unlock()

			c.Set("Retry-After", strconv.Itoa(r))
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "blocked",
				"message":     "request blocked",
				"retry_after": r,
			})
		}

		log.Printf("[RATE-ALLOW] untrusted-first-request ip=%s path=%s ua=%s fingerprint=%s",
			preferRealIP(c), c.Path(), c.Get("User-Agent", ""), key)

		return c.Next()
	}
}
