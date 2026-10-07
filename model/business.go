package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Business struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Phone         string             `bson:"phone" json:"phone"`
	Address       string             `bson:"address" json:"address"`
	TelegramToken string             `bson:"telegram_token" json:"telegram_token"`
	WorkingHours  map[string]string  `bson:"working_hours" json:"working_hours"`
	SlotDuration  int                `bson:"slot_duration" json:"slot_duration"` // به دقیقه
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

type BusinessList []Business
