package api

import (
	"Cidna/model"
	"Cidna/uscase"
	viewmodel "Cidna/viewModel"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResponse struct {
	Id        int    `json:"id"`
	Name_user string `json:"name"`
	Age       int    `json:"age"`
}

type UserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type ErrorRespons struct {
	Cod     int    `json:"cod"`
	Message string `json:"message"`
}

func GetUser(c *fiber.Ctx) error {

	id := c.QueryInt("id")
	name := c.Params("name")
	user := UserResponse{
		Id:        id,
		Name_user: name,
		Age:       20,
	}
	// خروجی متن
	// return c.SendString("hi")

	return c.JSON(user)
}

func AddUser(c *fiber.Ctx) error {

	user := new(viewmodel.GetUser)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(ErrorRespons{
			Cod:     104,
			Message: "Bad Request",
		})
	}

	// جایی که * استفاده کردیم برای مقدار دهی از & استفاده میکنیم
	u := &model.User{
		ID:      primitive.NewObjectID(),
		Name:    user.Name,
		AgeUser: user.AgeUser,
	}
	uscase.UserUse.SaveUser(u)

	return c.Status(201).JSON(ErrorRespons{
		Cod:     201,
		Message: "success save user",
	})
}
