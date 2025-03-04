package main

import (
	"fmt"
	"log"

	"github.com/ahnafasif/MarauderBoard/configs"
	"github.com/ahnafasif/MarauderBoard/controllers"
	"github.com/ahnafasif/MarauderBoard/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var app *fiber.App

func init() {
	configs.LoadEnv()
	database.ConnectDB()

	engine := html.New("./views", ".html")

	app = fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./public")

	controllers.RegisterRoutes(app)
}

func main() {
	port := configs.Port
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
