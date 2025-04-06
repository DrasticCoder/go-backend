package controllers

import (
	"encoding/json"
	"fmt"
	"go-backend/config"
	"go-backend/models"
	"go-backend/utils"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreatePost godoc
// @Summary Create a new post
// @Description Allows authors or admins to create posts (draft, scheduled, or published)
// @Tags Posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post body models.Post true "Post object"
// @Success 201 {object} models.Post
// @Failure 400 {string} string "Invalid input"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal error"
// @Router /api/v1/posts [post]
func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := utils.ExtractUserIDFromRequest(r)
fmt.Println("Extracted userID from JWT:", userID)

	role, _ := utils.ExtractUserRole(r)

	if role != "author" && role != "admin" {
		http.Error(w, "Only authors or admins can create posts", http.StatusForbidden)
		return
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	authorID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	post.ID = primitive.NewObjectID()
	post.AuthorID = authorID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.Status = "draft"

	if post.ScheduledAt != nil {
		if post.ScheduledAt.Before(time.Now()) {
			post.Status = "published"
			post.PublishedAt = post.ScheduledAt
		} else {
			post.Status = "scheduled"
		}
	}

	_, err = config.PostCollection.InsertOne(r.Context(), post)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// ListPosts godoc
// @Summary List published posts
// @Description Returns public posts, filtered by type and visibility
// @Tags Posts
// @Produce json
// @Param type query string false "Post type (idea, trade)"
// @Param visibility query string false "Post visibility (public, private, premium)"
// @Success 200 {array} models.Post
// @Failure 500 {string} string "Failed to fetch posts"
// @Router /api/v1/posts [get]
func ListPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	query := r.URL.Query()
	page := 1
	limit := 10
	if p := query.Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if l := query.Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}
	skip := (page - 1) * limit

	// Build filter
	filter := bson.M{}
	if status := query.Get("status"); status != "" {
		filter["status"] = strings.ToLower(status)
	} else {
		filter["status"] = "published" // Default to published if no status specified
	}

	// Type filter
	if postType := query.Get("type"); postType != "" {
		filter["type"] = strings.ToLower(postType)
	}

	// Visibility filter
	if visibility := query.Get("visibility"); visibility != "" {
		filter["visibility"] = strings.ToLower(visibility)
	}

	// Search
	if search := query.Get("search"); search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"content": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Sort
	sortField := "created_at"
	sortOrder := -1
	if sort := query.Get("sort"); sort != "" {
		if strings.HasPrefix(sort, "-") {
			sortField = strings.TrimPrefix(sort, "-")
			sortOrder = -1
		} else {
			sortField = sort
			sortOrder = 1
		}
	}

	// Build options
	findOptions := options.Find().
		SetSort(bson.D{{sortField, sortOrder}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// Get total count
	total, err := config.PostCollection.CountDocuments(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to count posts", http.StatusInternalServerError)
		return
	}

	// Execute query
	cursor, err := config.PostCollection.Find(r.Context(), filter, findOptions)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	posts := []models.Post{}
	if err := cursor.All(r.Context(), &posts); err != nil {
		http.Error(w, "Failed to parse posts", http.StatusInternalServerError)
		return
	}

	// Build response
	response := map[string]interface{}{
		"posts": posts,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": math.Ceil(float64(total) / float64(limit)),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// GetPost godoc
// @Summary Get a post by ID
// @Tags Posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} models.Post
// @Failure 404 {string} string "Post not found"
// @Failure 400 {string} string "Invalid ID"
// @Router /api/v1/posts/{id} [get]
func GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	postID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post models.Post
	err = config.PostCollection.FindOne(r.Context(), bson.M{"_id": postID}).Decode(&post)
	if err != nil || post.Status != "published" {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(post)
}

// UpdatePost godoc
// @Summary Update a post
// @Description Authors can update their own posts
// @Tags Posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body models.Post true "Post object"
// @Success 200 {string} string "Post updated"
// @Failure 400 {string} string "Bad input"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Post not found"
// @Router /api/v1/posts/{id} [put]
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := utils.ExtractUserIDFromRequest(r)
	params := mux.Vars(r)
	postID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var updates models.Post
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updates.UpdatedAt = time.Now()
	authorID, _ := primitive.ObjectIDFromHex(userID)

	filter := bson.M{"_id": postID, "author_id": authorID}
	update := bson.M{"$set": updates}

	result, err := config.PostCollection.UpdateOne(r.Context(), filter, update)
	if err != nil || result.MatchedCount == 0 {
		http.Error(w, "Update failed or not authorized", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Post updated"})
}

// DeletePost godoc
// @Summary Delete a post
// @Description Authors can delete their own posts
// @Tags Posts
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {string} string "Post deleted"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /api/v1/posts/{id} [delete]
func DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := utils.ExtractUserIDFromRequest(r)
	params := mux.Vars(r)
	postID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	authorID, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": postID, "author_id": authorID}

	result, err := config.PostCollection.DeleteOne(r.Context(), filter)
	if err != nil || result.DeletedCount == 0 {
		http.Error(w, "Delete failed or not authorized", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted"})
}

// ListMyPosts godoc
// @Summary List posts by current user
// @Description Returns all posts by the logged-in user
// @Tags Posts
// @Security BearerAuth
// @Success 200 {array} models.Post
// @Failure 500 {string} string "Server error"
// @Router /api/v1/posts/me [get]
func ListMyPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, _ := utils.ExtractUserIDFromRequest(r)
	authorID, _ := primitive.ObjectIDFromHex(userID)

	filter := bson.M{"author_id": authorID}

	cursor, err := config.PostCollection.Find(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	posts := []models.Post{}
	if err := cursor.All(r.Context(), &posts); err != nil {
		http.Error(w, "Parse error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}
