package notifications_controller

import (
	"github.com/ahnafasif/MarauderBoard/controllers/auth"
	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/models"
	load_locals "github.com/ahnafasif/MarauderBoard/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterNotificationController(app fiber.Router) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		data := load_locals.LoadLocals(ctx)

		data["PageName"] = "Notifications"
		return ctx.Render("notifications/index", data, "layouts/dashboard")
	})

	app.Get("/test", func(ctx *fiber.Ctx) error {
		user := ctx.Locals("User").(auth.UserSessionData)
		test_notification := &models.Notification{
			Title:  "Test Notification",
			Body:   "This is a test notification",
			UserId: user.ID,
			Seen:   false,
			Link:   nil,
		}
		_, err := models.AddNotification(database.DB, test_notification)
		if err != nil {
			return ctx.SendString("Error hoise")
		}
		return ctx.SendString("Done")
	})
}
