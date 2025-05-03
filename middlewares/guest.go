package middlewares

import (
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/gofiber/fiber/v2"
)

func GuestMiddleware(ctx *fiber.Ctx) error {
	_, err := helpers.GetAuthUserSessionData(ctx)
	if err == nil {
		return ctx.Redirect("/", fiber.StatusFound)
	}
	return ctx.Next()
}
