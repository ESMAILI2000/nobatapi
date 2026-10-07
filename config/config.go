package config

import "os"

type Config struct {
	TelegramToken string
	MongoURI      string
	DBName        string
}

func LoadConfig() *Config {
	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		MongoURI:      os.Getenv("MONGO_URI"),
		DBName:        os.Getenv("DB_NAME"),
	}
}
