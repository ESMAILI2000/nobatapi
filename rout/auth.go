package rout

import (
	"NobatApi/handlers"

	"github.com/gofiber/fiber/v2"
)

func AuthRouters(app *fiber.App) {
	app.Get("/login", handlers.LoginPage)
	app.Post("/login", handlers.Login)
}
