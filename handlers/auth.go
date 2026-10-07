package handlers

import (
	"NobatApi/uscase"

	"github.com/gofiber/fiber/v2"
)

// GET /login
func LoginPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

// POST /login
func Login(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// بررسی خالی نبودن ورودی‌ها
	if username == "" || password == "" {
		return c.Render("login", fiber.Map{
			"Error": "نام کاربری و رمز عبور را وارد کنید",
		})
	}

	// جستجوی کاربر در دیتابیس
	user, err := uscase.UserUse.GetUserByUserName(username)
	if err != nil {
		return c.Render("login", fiber.Map{
			"Error": "نام کاربری یا رمز عبور اشتباه است",
		})
	}

	// بررسی رمز عبور
	if user.Password != password {
		return c.Render("login", fiber.Map{
			"Error": "نام کاربری یا رمز عبور اشتباه است",
		})
	}

	// ورود موفق - انتقال به داشبورد
	return c.Redirect("/dashboard")
}
