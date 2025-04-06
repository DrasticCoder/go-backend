package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"go-backend/config"
	"go-backend/models"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration details"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Bad Request"
// @Router /api/v1/auth/register [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user.Password = strings.TrimSpace(user.Password)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	existing := config.UserCollection.FindOne(context.TODO(), bson.M{"email": user.Email})
	if existing.Err() == nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error while hashing password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPwd)
	user.CreatedAt = time.Now()

	res, err := config.UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	user.Password = "" // do not return hashed password

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// LoginUser godoc
// @Summary Login a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body AuthPayload true "Login Credentials"
// @Success 200 {string} string "JWT Token"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/v1/auth/login [post]
func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload AuthPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user models.User
	result := config.UserCollection.FindOne(context.TODO(), bson.M{"email": payload.Email})
	if result.Err() != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := result.Decode(&user); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user": map[string]interface{}{
			"id":   user.ID.Hex(),
			"role": user.Role,
		},
	})
}

// LogoutUser godoc
// @Summary Logout a user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {string} string "Logged out successfully"
// @Router /api/v1/auth/logout [post]
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
