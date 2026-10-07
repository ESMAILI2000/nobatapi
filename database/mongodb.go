package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongo() (*mongo.Database, error) {

	// اتصال به mongo در آینده یوزر و پس را داخل لینک اضافه میکنیم
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	// قطع اتصال دیتابیس هنگام متوقف کردن برنامه
	// defer client.Disconnect(context.TODO())

	if err != nil {
		log.Println(err.Error())
		log.Fatal(err) // هنگام خطا برنامه متوقف میشود این خط حذف کنید
		return nil, err
	}

	db := client.Database("NobatDb")
	return db, nil
}
