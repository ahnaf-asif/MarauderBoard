package main

import (
	"fmt"
	"log"

	"github.com/ahnafasif/MarauderBoard/configs"
	"github.com/ahnafasif/MarauderBoard/controllers"
	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/ahnafasif/MarauderBoard/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/shareed2k/goth_fiber"
)

var app *fiber.App

func init() {
	configs.LoadEnv()
	database.ConnectDB()

	engine := html.New("./views", ".html")
	engine.Reload(true)

	app = fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./public")
	goth_fiber.SessionStore = helpers.SetSessionConfig()
	helpers.SetAuthProviders()
	controllers.RegisterRoutes(app)
}

func main() {
	port := configs.Port
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
