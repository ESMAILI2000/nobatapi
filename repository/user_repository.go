package repository

import (
	"Cidna/model"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Save(user *model.User) error
}

type UserRepositoryMongo struct {
	collection *mongo.Collection
}

func NewUserRepositoryMongo(collection *mongo.Collection) *UserRepositoryMongo {
	return &UserRepositoryMongo{
		collection: collection,
	}
}

func (r *UserRepositoryMongo) Save(user *model.User) error {
	_, err := r.collection.InsertOne(context.TODO(), user)
	return err
}
