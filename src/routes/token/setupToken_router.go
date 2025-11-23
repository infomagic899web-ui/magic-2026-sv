package tokens

import (
	"github.com/gofiber/fiber/v3"
)

func SetupTokenRouter(app fiber.Router) {
	api := app.Group("/token")
	api.Get("/csrf", GetCSRFToken)
	api.Get("/rsp", GetRSToken)
	api.Get("/nonce", GetNonceToken)
	api.Get("/sha512", GetSHA512Token)
	api.Get("/sha256", GetSHA256Token)
}
