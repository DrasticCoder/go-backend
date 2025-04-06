package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
	VisibilityPremium Visibility = "premium"
)

type PostType string

const (
	PostTypeIdea  PostType = "idea"
	PostTypeTrade PostType = "trade"
)

type Post struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AuthorID    primitive.ObjectID `bson:"author_id" json:"author_id"`
	Title       string             `bson:"title" json:"title"`
	Content     string             `bson:"content" json:"content"` // markdown
	Tags        []string           `bson:"tags" json:"tags"`
	Visibility  Visibility         `bson:"visibility" json:"visibility"` // public, private, premium
	Type        PostType           `bson:"type" json:"type"`             // idea, trade
	Status      string             `bson:"status" json:"status"`         // draft, published, scheduled
	MediaURLs   []string           `bson:"media_urls,omitempty" json:"media_urls,omitempty"`
	ScheduledAt *time.Time         `bson:"scheduled_at,omitempty" json:"scheduled_at,omitempty"`
	PublishedAt *time.Time         `bson:"published_at,omitempty" json:"published_at,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
