package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func StartPostScheduler(db *mongo.Database) {
	ticker := time.NewTicker(1 * time.Minute) // Runs every 1 min
	go func() {
		for range ticker.C {
			publishScheduledPosts(db)
		}
	}()
}

func publishScheduledPosts(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postsColl := db.Collection("posts")

	filter := bson.M{
		"published":    false,
		"scheduledFor": bson.M{"$lte": time.Now()},
	}

	update := bson.M{
		"$set": bson.M{"published": true, "publishedAt": time.Now()},
	}

	result, err := postsColl.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Printf("Failed to publish scheduled posts: %v", err)
		return
	}

	if result.ModifiedCount > 0 {
		log.Printf("✅ Published %d scheduled post(s)", result.ModifiedCount)
	}
}
