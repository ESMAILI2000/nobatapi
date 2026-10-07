package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Username   string             `bson:"username"`
	Password   string             `bson:"password"` // باید هش شود
	BusinessID primitive.ObjectID `bson:"business_id"`
}

type UserList []User
