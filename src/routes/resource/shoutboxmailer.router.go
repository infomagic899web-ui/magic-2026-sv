package resources

import (
	"magic899-server/src/controllers"
	"magic899-server/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func ShoutboxMailerRouter(router fiber.Router) {

	api := router.Group("/shoutbox-mailer")
	api.Post("/", middlewares.RateLimiterMiddleware(), middlewares.CSRFTokenMiddleware, controllers.SendAutoReplyShoutboxMailer)
}
