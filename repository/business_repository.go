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

type BusinessRepository interface {
	Save(business *model.Business) error
	Update(id primitive.ObjectID, business *model.Business) error
	FindByID(id primitive.ObjectID) (*model.Business, error)
	Delete(id primitive.ObjectID) error
	FindAllByFilter(name string, ownerID primitive.ObjectID, isActive *bool, pageNum int64, pageSize int64, sortBy string, isDescending bool) (*model.BusinessList, error)
}

type BusinessRepositoryMongo struct {
	collection *mongo.Collection
}

func NewBusinessRepositoryMongo(collection *mongo.Collection) *BusinessRepositoryMongo {
	return &BusinessRepositoryMongo{collection: collection}
}

// -------------------- CRUD METHODS --------------------

func (r *BusinessRepositoryMongo) Save(business *model.Business) error {
	_, err := r.collection.InsertOne(context.TODO(), business)
	return err
}

func (r *BusinessRepositoryMongo) Update(id primitive.ObjectID, business *model.Business) error {
	update := bson.M{"$set": business}
	_, err := r.collection.UpdateByID(context.TODO(), id, update)
	return err
}

func (r *BusinessRepositoryMongo) FindByID(id primitive.ObjectID) (*model.Business, error) {
	var result model.Business
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *BusinessRepositoryMongo) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

// -------------------- FIND ALL WITH FILTER (PANEL) --------------------
// فیلتر داینامیک بر اساس نام بیزینس، مالک (owner_id) و وضعیت فعال بودن

func (r *BusinessRepositoryMongo) FindAllByFilter(
	name string,
	ownerID primitive.ObjectID,
	isActive *bool,
	pageNum int64,
	pageSize int64,
	sortBy string,
	isDescending bool,
) (*model.BusinessList, error) {

	var opt options.FindOptions

	// تنظیم مرتب‌سازی
	if sortBy != "" {
		if isDescending {
			opt.SetSort(bson.D{{Key: sortBy, Value: -1}})
		} else {
			opt.SetSort(bson.D{{Key: sortBy, Value: 1}})
		}
	} else {
		sortBy = "created_at"
		opt.SetSort(bson.D{{Key: sortBy, Value: -1}})
	}

	// تنظیم صفحه‌بندی
	opt.SetSkip(pageNum * pageSize)
	opt.SetLimit(pageSize)

	// ساخت کوئری پویا
	andQuery := []bson.M{}

	if name != "" {
		andQuery = append(andQuery, bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: name}}})
	}

	if ownerID != primitive.NilObjectID {
		andQuery = append(andQuery, bson.M{"owner_id": ownerID})
	}

	if isActive != nil {
		andQuery = append(andQuery, bson.M{"is_active": *isActive})
	}

	// در صورت خالی بودن فیلترها
	if len(andQuery) == 0 {
		andQuery = append(andQuery, bson.M{})
	}

	filter := bson.M{"$and": andQuery}

	cursor, err := r.collection.Find(context.TODO(), filter, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	allBusiness := make(model.BusinessList, 0)
	for cursor.Next(context.TODO()) {
		var result model.Business
		if err := cursor.Decode(&result); err != nil {
			log.Println("decode error:", err)
			continue
		}
		allBusiness = append(allBusiness, result)
	}

	return &allBusiness, nil
}
