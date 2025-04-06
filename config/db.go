package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo *mongo.Client
var UserCollection *mongo.Collection
var PostCollection *mongo.Collection

func InitDB() {
    uri := os.Getenv("MONGO_URI")
    if uri == "" {
        uri = "mongodb://localhost:27017"
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatalf("❌ MongoDB connection error: %v", err)
    }

    Mongo = client
    UserCollection = Mongo.Database("crm").Collection("users")
    PostCollection = Mongo.Database("crm").Collection("posts")
    Logger.Info("📦 Connected to MongoDB!")
}
