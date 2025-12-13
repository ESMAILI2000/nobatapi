package main

import (
	"Cidna/database"
	"Cidna/rout"
	"Cidna/uscase"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// اتصال به دیتابیس
	db, err := database.ConnectToMongo()
	if err != nil {
		fmt.Println(err)
	}

	uscase.UscaseInit(db)

	app := fiber.New()

	rout.CidnaRouters(app)

	// app.Get("/:name", api.GetUser)
	// app.Post("/user", api.AddUser)

	app.Listen(":8081")

}
