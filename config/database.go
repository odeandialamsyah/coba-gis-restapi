package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb+srv://odeandi:13210911@cluster0.big5xbd.mongodb.net/"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "geo_app"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Mongo ping error:", err)
	}

	DB = client.Database(dbName)
	log.Println("✅ Connected to MongoDB:", dbName)
}
