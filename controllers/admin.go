package controllers

import (
	"encoding/json"
	"go-backend/config"
	"go-backend/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// ListUsers godoc
// @Summary List all users
// @Description Returns a list of all users (admin only) with pagination, search, sort and filters
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Param search query string false "Search by email or name"
// @Param sort query string false "Sort field (prefix with - for desc)"
// @Param role query string false "Filter by role"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {string} string "Unauthorized"
// @Router /api/v1/admin/users [get]
func ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	query := r.URL.Query()
	
	// Pagination
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit

	// Build filter
	filter := bson.M{}

	// Search
	if search := query.Get("search"); search != "" {
		filter["$or"] = []bson.M{
			{"email": bson.M{"$regex": search, "$options": "i"}},
			{"name": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Role filter
	if role := query.Get("role"); role != "" {
		filter["role"] = role
	}

	// Sorting
	sort := bson.D{}
	if sortParam := query.Get("sort"); sortParam != "" {
		sortFields := strings.Split(sortParam, ",")
		for _, field := range sortFields {
			if strings.HasPrefix(field, "-") {
				sort = append(sort, bson.E{Key: strings.TrimPrefix(field, "-"), Value: -1})
			} else {
				sort = append(sort, bson.E{Key: field, Value: 1})
			}
		}
	}

	// Set options
	findOptions := options.Find()
	findOptions.SetSort(sort)
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	// Get total count
	total, err := config.UserCollection.CountDocuments(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get users
	users := []models.User{}
	cursor, err := config.UserCollection.Find(r.Context(), filter, findOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	if err := cursor.All(r.Context(), &users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// GetUser godoc
// @Summary Get single user
// @Description Returns a single user by ID (admin only)
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Router /api/v1/admin/users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	err = config.UserCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user (admin only)
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 201 {object} models.User
// @Failure 401 {string} string "Unauthorized"
// @Failure 400 {string} string "Invalid request body"
// @Router /api/v1/admin/users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert user
	result, err := config.UserCollection.InsertOne(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user role or email (admin only)
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "User object"
// @Success 200 {object} models.User
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Router /api/v1/admin/users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updates models.User
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build update document with only provided fields
	updateFields := bson.M{}
	if updates.Email != "" {
		updateFields["email"] = updates.Email
	}
	if updates.Role != "" {
		updateFields["role"] = updates.Role
	}

	result, err := config.UserCollection.UpdateOne(
		r.Context(),
		bson.M{"_id": id},
		bson.M{"$set": updateFields},
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updatedUser models.User
	config.UserCollection.FindOne(r.Context(), bson.M{"_id": id}).Decode(&updatedUser)
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user (admin only)
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Router /api/v1/admin/users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	result, err := config.UserCollection.DeleteOne(r.Context(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
