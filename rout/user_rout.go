package rout

import (
	"Cidna/api"

	"github.com/gofiber/fiber/v2"
)

func UserRouters(app *fiber.App) {
	app.Get("/:name", api.GetUser)
	app.Post("/user", api.AddUser)
}
