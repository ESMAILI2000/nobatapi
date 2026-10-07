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

type AppointmentRepository interface {
	Save(appointment *model.Appointment) error
	Update(id primitive.ObjectID, appointment *model.Appointment) error
	IsSlotTaken(date, timeSlot string) (bool, error)
	UpdateStatus(id primitive.ObjectID, status string) error
	FindByID(id primitive.ObjectID) (*model.Appointment, error)
	FindByDate(date string) ([]model.Appointment, error)
	Delete(id primitive.ObjectID) error
	FindAll() ([]model.Appointment, error)
	FindAllAppointmentByBusinessID(id primitive.ObjectID) (*model.AppointmentList, error)
	FindAllByFilter(
		businessId primitive.ObjectID, serviceId primitive.ObjectID, clientId primitive.ObjectID, staffId primitive.ObjectID,
		status string, startDate *primitive.DateTime, endDate *primitive.DateTime,
		pageNum int64, pageSize int64, sortBy string, isDescending bool,
	) (*model.AppointmentList, error)
}

type AppointmentRepositoryMongo struct {
	collection *mongo.Collection
}

func NewAppointmentRepositoryMongo(collection *mongo.Collection) *AppointmentRepositoryMongo {
	return &AppointmentRepositoryMongo{collection: collection}
}

// -------------------- CRUD METHODS --------------------

func (r *AppointmentRepositoryMongo) Save(appointment *model.Appointment) error {
	_, err := r.collection.InsertOne(context.TODO(), appointment)
	return err
}

func (r *AppointmentRepositoryMongo) Update(id primitive.ObjectID, appointment *model.Appointment) error {
	update := bson.M{"$set": appointment}
	_, err := r.collection.UpdateByID(context.TODO(), id, update)
	return err
}

// چک می‌کند آیا این ساعت برای این تاریخ قبلاً رزرو شده یا نه
func (r *AppointmentRepositoryMongo) IsSlotTaken(date, timeSlot string) (bool, error) {
	count, err := r.collection.CountDocuments(context.TODO(), bson.M{
		"date":   date,
		"time":   timeSlot,
		"status": bson.M{"$ne": "cancelled"},
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AppointmentRepositoryMongo) UpdateStatus(id primitive.ObjectID, status string) error {
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.collection.UpdateByID(context.TODO(), id, update)
	return err
}

func (r *AppointmentRepositoryMongo) FindByID(id primitive.ObjectID) (*model.Appointment, error) {
	var result model.Appointment

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *AppointmentRepositoryMongo) FindByDate(date string) ([]model.Appointment, error) {
	var results []model.Appointment
	cursor, err := r.collection.Find(context.TODO(), bson.M{"date": date})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var appt model.Appointment
		cursor.Decode(&appt)
		results = append(results, appt)
	}

	return results, nil
}

func (r *AppointmentRepositoryMongo) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func (r *AppointmentRepositoryMongo) FindAll() ([]model.Appointment, error) {
	var results []model.Appointment

	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var appt model.Appointment
		if err := cursor.Decode(&appt); err != nil {
			return nil, err
		}
		results = append(results, appt)
	}

	return results, nil
}

// همه نوبت های یک بیزینس
func (r *AppointmentRepositoryMongo) FindAllAppointmentByBusinessID(id primitive.ObjectID) (*model.AppointmentList, error) {
	var opt options.FindOptions
	var cursor *mongo.Cursor
	var err error

	cursor, err = r.collection.Find(context.TODO(), bson.M{"business_id": id}, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())
	allResult := make(model.AppointmentList, 0)
	for cursor.Next(context.TODO()) {

		var result model.Appointment
		if err = cursor.Decode(&result); err != nil {
			log.Println(err)
		}

		allResult = append(allResult, result)
	}

	return &allResult, nil
}

func (r *AppointmentRepositoryMongo) FindAllByFilter(
	businessId primitive.ObjectID,
	serviceId primitive.ObjectID,
	clientId primitive.ObjectID,
	staffId primitive.ObjectID,
	status string,
	startDate *primitive.DateTime, // از تاریخ
	endDate *primitive.DateTime, // تا تاریخ
	pageNum int64,
	pageSize int64,
	sortBy string,
	isDescending bool,
) (*model.AppointmentList, error) {

	// تنظیم گزینه‌های جستجو (صفحه‌بندی و مرتب‌سازی)
	var opt options.FindOptions

	// تنظیم مرتب‌سازی ( sortBy پیش‌فرض 'created_at' است اگر چیزی مشخص نشده باشد)
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := 1
	if isDescending {
		sortOrder = -1
	}
	opt.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})

	// تنظیم صفحه‌بندی
	opt.SetSkip(pageNum * pageSize)
	opt.SetLimit(pageSize)

	// ساخت کوئری پویا با استفاده از $and
	andQuery := []bson.M{}

	if businessId != primitive.NilObjectID {
		andQuery = append(andQuery, bson.M{"business_id": businessId})
	}
	if serviceId != primitive.NilObjectID {
		andQuery = append(andQuery, bson.M{"service_id": serviceId})
	}
	if clientId != primitive.NilObjectID {
		andQuery = append(andQuery, bson.M{"client_id": clientId})
	}
	if staffId != primitive.NilObjectID {
		andQuery = append(andQuery, bson.M{"staff_id": staffId})
	}
	if status != "" {
		andQuery = append(andQuery, bson.M{"status": status})
	}
	if startDate != nil {
		andQuery = append(andQuery, bson.M{
			"date": bson.M{"$gte": *startDate},
		})
	}
	if endDate != nil {
		andQuery = append(andQuery, bson.M{
			"date": bson.M{"$lte": *endDate},
		})
	}

	var filter bson.M
	if len(andQuery) > 0 {
		filter = bson.M{"$and": andQuery}
	} else {
		filter = bson.M{} // یا new(bson.M)
	}

	// اجرای جستجو
	cursor, err := r.collection.Find(context.TODO(), filter, &opt)
	if err != nil {
		log.Printf("Error finding appointments: %v\n", err)
		return nil, err
	}

	defer cursor.Close(context.TODO())

	allAppointments := make(model.AppointmentList, 0)
	// پیمایش نتایج و دیکد کردن هر سند
	for cursor.Next(context.TODO()) {
		var result model.Appointment
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Error decoding appointment: %v\n", err)
			// در صورت خطای دیکد، ممکن است بخواهیم حلقه را بشکنیم یا فقط خطا را لاگ کنیم
			continue // ادامه با سند بعدی
		}
		allAppointments = append(allAppointments, result)
	}

	// بررسی خطای احتمالی در طول پیمایش cursor
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v\n", err)
		return nil, err
	}

	return &allAppointments, nil
}
