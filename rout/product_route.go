package rout

import (
	"Cidna/api"

	"github.com/gofiber/fiber/v2"
)

func ProductRouters(app *fiber.App) {
	app.Get("/:name", api.GetUser)
	app.Post("/user", api.AddUser)

}
