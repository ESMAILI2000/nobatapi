package uscase

import (
	"Cidna/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

func UscaseInit(db *mongo.Database) {
	UserUse = NewUserUsecase(repository.NewUserRepositoryMongo(db.Collection("user")))
}
