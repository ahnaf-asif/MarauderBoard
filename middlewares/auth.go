package middlewares

import (
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	user, err := helpers.GetAuthUserSessionData(ctx)
	if err != nil {
		ctx.Cookie(&fiber.Cookie{
			Name:     "redirect_uri",
			MaxAge:   300,
			Value:    ctx.OriginalURL(),
			HTTPOnly: true,
			SameSite: "Laz",
		})

		return ctx.Redirect("/auth/google")
	}

	ctx.Locals("User", user)

	return ctx.Next()
}
