package middlewares

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

var (
	requests   = make(map[string]int)
	timestamps = make(map[string]time.Time)
	mu         sync.Mutex
)

func RateLimiterMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()
		now := time.Now()

		mu.Lock()
		defer mu.Unlock()

		// if no record, init
		if _, ok := timestamps[ip]; !ok {
			timestamps[ip] = now
			requests[ip] = 0
		}

		// check if current window expired
		if now.Sub(timestamps[ip]) > time.Minute {
			// reset for a new 1-minute window
			requests[ip] = 0
			timestamps[ip] = now
		}

		// increment request count
		requests[ip]++

		if requests[ip] > 5 {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, slow down",
			})
		}

		return c.Next()
	}
}
