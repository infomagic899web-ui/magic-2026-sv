package resources

import (
	"magic-server-2026/src/controllers"
	"magic-server-2026/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func ShoutboxMailerRouter(router fiber.Router) {

	api := router.Group("/shoutbox-mailer")
	api.Post("/", middlewares.RateLimiterMiddleware(), middlewares.CSRFTokenMiddleware, controllers.SendAutoReplyShoutboxMailer)
}
