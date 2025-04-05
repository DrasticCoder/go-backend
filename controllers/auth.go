package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
// @Param user body models.User true "User Info (default: {username: 'john123', email: 'john@doe.com', password: 'Pa$$w0rd!', role: 'free'})"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Bad Request"
// @Router /api/v1/auth/register [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user.Password = strings.TrimSpace(user.Password)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPwd)
	user.CreatedAt = time.Now()

	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "User already exists or invalid data", http.StatusBadRequest)
		return
	}
	// user.Password = ""
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
	var payload AuthPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user models.User
	result := config.DB.Where("email = ?", payload.Email).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	

	// Compare password
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
	// 	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	// 	return
	// }

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// LogoutUser godoc
// @Summary Logout a user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {string} string "Logged out successfully"
// @Router /api/v1/auth/logout [post]
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// For a stateless JWT implementation, client-side token removal is sufficient
	// Server could implement a blacklist/revocation mechanism for additional security
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}


