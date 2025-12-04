package api

import "github.com/gofiber/fiber/v2"

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

	user := new(UserRequest)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(ErrorRespons{
			Cod:     104,
			Message: "Bad Request",
		})
	}
	// خروجی متن
	// return c.SendString("hi")

	return c.JSON(UserResponse{
		Id:        1,
		Name_user: user.Name,
		Age:       user.Age,
	})
}
