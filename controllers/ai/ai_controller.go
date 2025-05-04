package ai_controller

import (
	"log"

	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/ahnafasif/MarauderBoard/services"
	load_locals "github.com/ahnafasif/MarauderBoard/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	model = "gemma:2b"
)

func RegisterAiControllers(app fiber.Router) {
	app.Post("/ask", func(ctx *fiber.Ctx) error {
		question := ctx.FormValue("question")
		if question == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing question",
			})
		}
		userMessage := services.Message{
			Role:    "user",
			Content: question,
		}
		log.Println("New Question: ", userMessage.Content)
		messages := []services.Message{userMessage}

		responseMessage, err := services.SendToOllama(model, messages)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get response from Ollama",
			})
		}
		log.Println("New Response: ", responseMessage.Content)
		tmp := responseMessage.Content
		responseMessage.Content, err = helpers.MarkdownWithMath(responseMessage.Content)
		if err != nil {
			log.Println("Error rendering markdown: ", err)
			responseMessage.Content = tmp
		} else {
			log.Println("Rendered Markdown: ", responseMessage.Content)
		}
		renderedHTML, err := services.RenderMessagesHTML([]services.Message{responseMessage})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to render message HTML",
			})
		}

		log.Println("Rendered HTML: ", renderedHTML)

		return ctx.SendString(renderedHTML)
	})

	app.Get("/chat", func(ctx *fiber.Ctx) error {
		session_id := services.GetOrCreateOllamaSession(ctx)
		messages := services.GetMessages(session_id)

		data := load_locals.LoadLocals(ctx)
		data["PageTitle"] = "AI Chat"
		data["Messages"] = messages

		return ctx.Render("ai/chat", data, "layouts/dashboard")
	})
}
