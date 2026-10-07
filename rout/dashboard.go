package rout

import (
	"NobatApi/handlers"

	"github.com/gofiber/fiber/v2"
)

func DashboardRouters(app *fiber.App) {
	app.Get("/dashboard", handlers.Dashboard)
	app.Post("/appointment/create", handlers.CreateAppointment)
	app.Post("/appointment/update-status", handlers.UpdateAppointmentStatus)
}
