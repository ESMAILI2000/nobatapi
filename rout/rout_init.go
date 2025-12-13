package rout

import "github.com/gofiber/fiber/v2"

func CidnaRouters(app *fiber.App) {
	UserRouters(app)
	ProductRouters(app)
}
