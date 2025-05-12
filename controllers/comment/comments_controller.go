package comments_controller

import (
	"log"
	"strconv"

	"github.com/ahnafasif/MarauderBoard/controllers/auth"
	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/models"
	load_locals "github.com/ahnafasif/MarauderBoard/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterCommentsController(app fiber.Router) {
	app.Post("/new", func(ctx *fiber.Ctx) error {
		task_id := ctx.Params("task_id")
		task_id_int, _ := strconv.Atoi(task_id)

		content := ctx.FormValue("content")
		data := load_locals.LoadLocals(ctx)
		user := data["User"].(auth.UserSessionData)

		log.Println("User ID:", user.ID)
		log.Println("Task ID:", task_id_int)
		log.Println("Content:", content)
		log.Println("Content", content)

		comment := &models.Comment{
			Content: content,
			TaskId:  uint(task_id_int),
			UserId:  user.ID,
		}

		db_comment, _ := models.AddNewComment(database.DB, comment)
		return ctx.Render("commentItem", db_comment)
	})

	app.Post("/reply/:comment_id", func(ctx *fiber.Ctx) error {
		task_id := ctx.Params("task_id")
		task_id_int, _ := strconv.Atoi(task_id)
		comment_id := ctx.Params("comment_id")
		comment_id_int, _ := strconv.Atoi(comment_id)
		comment_id_uint := uint(comment_id_int)
		content := ctx.FormValue("content")
		data := load_locals.LoadLocals(ctx)
		user := data["User"].(auth.UserSessionData)

		comment := &models.Comment{
			Content:  content,
			TaskId:   uint(task_id_int),
			UserId:   user.ID,
			ParentId: &comment_id_uint,
		}
		db_comment, _ := models.AddNewComment(database.DB, comment)
		return ctx.Render("commentItem", db_comment)
	})
}
