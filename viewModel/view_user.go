package viewmodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetUser struct {
	Name    string `json:"name"`
	AgeUser int    `json:"age_user"`
}

type SetUser struct {
	Id      primitive.ObjectID `json:"id"`
	Name    string             `json:"name"`
	AgeUser int                `json:"age_user"`
}

type SetUserList []SetUser
