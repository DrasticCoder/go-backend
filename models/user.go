package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name  string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    Password  string             `bson:"password" json:"password"`
    Role      string             `bson:"role" json:"role"`
    IsActive  bool               `bson:"is_active" json:"is_active"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
