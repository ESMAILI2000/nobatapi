package handlers

import (
	"NobatApi/model"
	"NobatApi/uscase"

	"github.com/gofiber/fiber/v2"
)

type ErrorRespons struct {
	Cod     int    `json:"cod"`
	Message string `json:"message"`
}

// ثبت کاربر
func AddUser(c *fiber.Ctx) error {

	u := &model.User{
		Username: "admin",
		Password: "123456",
	}
	err := uscase.UserUse.SaveUser(u) // حتما خطا را دریافت کن
	if err != nil {

		return c.Status(500).JSON(ErrorRespons{
			Cod:     500,
			Message: "خطا در ذخیره سازی ",
		})

	}

	return c.Status(201).JSON(ErrorRespons{
		Cod:     201,
		Message: "success save user",
	})
}
