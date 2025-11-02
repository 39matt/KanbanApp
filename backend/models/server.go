package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var db *mongo.Database

func InitServer(ctx context.Context) *mongo.Database {
	databaseURI := os.Getenv("MONGO_URI")
	if databaseURI == "" {
		log.Fatal("MONGO_URI not set in environment")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().
		ApplyURI(databaseURI).
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB")

	db = client.Database(os.Getenv("MONGO_DB"))
	return db
}
