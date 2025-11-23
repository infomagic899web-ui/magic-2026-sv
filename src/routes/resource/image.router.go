package resources

import (
	"magic-server-2026/src/controllers"
	"magic-server-2026/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func ImageRouter(router fiber.Router) {

	api := router.Group("/view")
	api.Post("/", middlewares.CSRFTokenMiddleware, controllers.UploadImageHandler)
	api.Get("/:filename", controllers.GetImageHandler)

}
