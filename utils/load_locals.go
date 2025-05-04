package load_locals

import "github.com/gofiber/fiber/v2"

func LoadLocals(ctx *fiber.Ctx) fiber.Map {
	user := ctx.Locals("User")
	unseen_notifications := ctx.Locals("UnseenNotifications")
	total_unseen_notifications := ctx.Locals("UnseenCount")

	mp := fiber.Map{
		"User":                user,
		"UnseenCount":         total_unseen_notifications,
		"UnseenNotifications": unseen_notifications,
	}

	return mp
}
