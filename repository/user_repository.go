package repository

import (
	"NobatApi/model"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	Save(user *model.User) error
	Update(id primitive.ObjectID, user *model.User) error
	FindByID(id primitive.ObjectID) (*model.User, error)
	FindByUsername(userName string) (*model.User, error)
	Delete(id primitive.ObjectID) error
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

func (r *UserRepositoryMongo) Update(id primitive.ObjectID, user *model.User) error {
	update := bson.M{
		"$set": user,
	}
	_, err := r.collection.UpdateByID(context.TODO(), id, update)
	return err
}

func (r *UserRepositoryMongo) FindByID(id primitive.ObjectID) (*model.User, error) {

	var result model.User

	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *UserRepositoryMongo) FindByUsername(userName string) (*model.User, error) {

	var result model.User

	err := r.collection.FindOne(context.TODO(), bson.M{"username": userName}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *UserRepositoryMongo) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func (r *UserRepositoryMongo) FindAllByFilterPanel(name string, role string, pageNum int64, pageSize int64, sortBy string, isDescending bool) (*model.UserList, error) {

	var opt options.FindOptions

	if sortBy != "" {
		if isDescending {
			opt.SetSort(bson.D{primitive.E{Key: sortBy, Value: -1}})
		} else {
			opt.SetSort(bson.D{primitive.E{Key: sortBy, Value: 1}})
		}
	} else {
		sortBy = "name"
	}
	opt.SetSkip(pageSize * pageNum)
	opt.SetLimit(pageSize)

	AndQuery := []bson.M{}

	if name != "" {
		a := bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: name}}}
		AndQuery = append(AndQuery, a)
	}

	if role != "" {
		b := bson.M{"role": role}
		AndQuery = append(AndQuery, b)
	}

	cursor, err := r.collection.Find(context.TODO(), bson.M{"$and": AndQuery}, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())
	allUser := make(model.UserList, 0)
	for cursor.Next(context.TODO()) {

		var result model.User
		if err = cursor.Decode(&result); err != nil {
			log.Println(err)
		}

		allUser = append(allUser, result)
	}
	return &allUser, nil
}
