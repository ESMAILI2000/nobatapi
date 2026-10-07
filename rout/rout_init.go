package rout

import "github.com/gofiber/fiber/v2"

func NobatRouters(app *fiber.App) {
	UserRouters(app)
	AuthRouters(app)
	DashboardRouters(app)
	DebugRouters(app)
}
