package middlewares

import (
	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/models"
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

		return ctx.Redirect("/auth/login")
	}

	unseen_notifications, err := models.GetUnseenNotificationsByUserId(database.DB, user.ID)
	if err != nil {
		unseen_notifications = []*models.Notification{}
	}

	total_unseen_notifications := len(unseen_notifications)

	ctx.Locals("User", user)
	ctx.Locals("UnseenCount", total_unseen_notifications)
	ctx.Locals("UnseenNotifications", unseen_notifications)

	return ctx.Next()
}
