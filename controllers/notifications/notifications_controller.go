package notifications_controller

import (
	"math/rand"
	"strconv"

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

	app.Post("/:id/mark-read", func(ctx *fiber.Ctx) error {
		notification_id, _ := ctx.ParamsInt("id")
		err := models.MarkNotificationAsSeen(database.DB, uint(notification_id))
		if err != nil {
			return ctx.Status(500).SendString("Error hoise")
		}
		return ctx.SendString("")
	})

	app.Get("/test", func(ctx *fiber.Ctx) error {
		user := ctx.Locals("User").(auth.UserSessionData)

		random_number := rand.Intn(100)
		title := "Test Notification " + strconv.Itoa(random_number)

		test_notification := &models.Notification{
			Title:  title,
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
