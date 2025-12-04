package main

import (
	"Cidna/api"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	app.Get("/:name", api.GetUser)
	app.Post("/user", api.AddUser)

	app.Listen(":8081")

}
