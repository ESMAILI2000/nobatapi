package uscase

import (
	"NobatApi/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

func UscaseInit(db *mongo.Database) {
	UserUse = NewUserUsecase(repository.NewUserRepositoryMongo(db.Collection("user")))
	AppointmentUse = NewAppointmentUsecase(repository.NewAppointmentRepositoryMongo(db.Collection("appointments")))
}
