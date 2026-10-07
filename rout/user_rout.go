package rout

import (
	"NobatApi/handlers"

	"github.com/gofiber/fiber/v2"
)

func UserRouters(app *fiber.App) {
	app.Post("/user", handlers.AddUser)
}
