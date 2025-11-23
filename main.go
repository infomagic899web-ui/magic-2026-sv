package main

import (
	"magic-server-2026/src/db"
	"magic-server-2026/src/gen"
	"magic-server-2026/src/middlewares"
	"magic-server-2026/src/routes"
	"magic-server-2026/src/server"
	"magic-server-2026/src/utils"

	"github.com/gofiber/fiber/v3"
)

func main() {
	utils.LoadEnv()
	db.Init()
	app := fiber.New(fiber.Config{
		EnableIPValidation: true,
		TrustProxy:         true,
	})

	middlewares.Setup(app)

	gen.Init()

	routes.SetupRouter(app)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World! You've just seeing me in production. :)")
	})

	server.Start(app)
}
