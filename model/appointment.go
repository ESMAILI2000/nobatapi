package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BusinessID     primitive.ObjectID `bson:"business_id" json:"business_id"`
	CustomerName   string             `bson:"customer_name" json:"customer_name"`
	CustomerChatID int64              `bson:"customer_chat_id" json:"customer_chat_id"`
	Date           string             `bson:"date" json:"date"`
	Time           string             `bson:"time" json:"time"`
	Status         string             `bson:"status" json:"status"`
	Source         string             `bson:"source" json:"source"` // "telegram", "walkin", "phone"
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
}

type AppointmentList []Appointment
