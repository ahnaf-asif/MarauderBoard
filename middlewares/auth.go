package middlewares

import (
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	user, err := helpers.GetAuthUserSessionData(ctx)
	if err != nil {
		return ctx.Redirect("/auth/google")
	}

	ctx.Locals("User", user)

	return ctx.Next()
}
