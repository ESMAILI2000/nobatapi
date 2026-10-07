package main

import (
	"NobatApi/bot"
	"NobatApi/config"
	"NobatApi/database"
	"NobatApi/rout"
	"NobatApi/uscase"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.LoadConfig()

	db, err := database.ConnectToMongo()
	if err != nil {
		fmt.Println(err)
	}

	// تنظیم موتور قالب‌ها
	engine := html.New("./templates", ".html")

	uscase.UscaseInit(db)

	// ساخت سرویس ربات
	botService := bot.NewBotService(cfg)

	// اجرای ربات در گوروتین جداگانه
	go botService.Start()

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// مسیرهای ایستا
	app.Static("/static", "./static")

	// ثبت همه روترها
	rout.NobatRouters(app)

	log.Fatal(app.Listen(":3000"))
}
