package routes

import (
	timein "magic-server-2026/src/TimeIn"
	"magic-server-2026/src/middlewares"
	resources "magic-server-2026/src/routes/resource"
	tokens "magic-server-2026/src/routes/token"

	"github.com/gofiber/fiber/v3"
)

func SetupRouter(app *fiber.App) {

	years := timein.Years(2)

	rl := middlewares.NewRateLimiter(years)

	api := app.Group("/api", middlewares.ReferrerMiddleware, middlewares.ResourceTokenMiddleware, middlewares.MakeTrustedRateLimiterMiddleware(rl))

	apiSecure := api.Group("/v1")

	resourceRoutes := []func(router fiber.Router){
		resources.NewsRouterResource,
		resources.MusicRouter,
		resources.ShowRouter,
		resources.YoutubeMagicVideosRouter,
		resources.MagicVideoRouter,
		resources.MoviesRouter,
		resources.RequestedSongRouter,
		resources.ShoutboxMailerRouter,
	}

	for _, r := range resourceRoutes {
		r(apiSecure)
	}

	resources.TestRouter(api)
	resources.ImageRouter(api)
	resources.GetPlayerRouter(api)
	tokens.SetupTokenRouter(api)
}
