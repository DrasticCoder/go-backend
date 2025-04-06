package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func ExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no token provided")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid token format")
	}

	return parts[1], nil
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid claims")
}

func ExtractUserIDFromRequest(r *http.Request) (string, error) {
	tokenString, err := ExtractToken(r)
	if err != nil {
		return "", err
	}

	claims, err := ParseJWT(tokenString)
	if err != nil {
		return "", err
	}

	if userID, ok := claims["user_id"].(string); ok {
		fmt.Println("userID", userID)
		return userID, nil
	}

	return "", errors.New("user ID not found in token")
}

// You can use this in middleware if needed
func ExtractUserRole(r *http.Request) (string, error) {
	tokenString, err := ExtractToken(r)
	if err != nil {
		return "", err
	}

	claims, err := ParseJWT(tokenString)
	if err != nil {
		return "", err
	}

	if role, ok := claims["role"].(string); ok {
		return role, nil
	}

	return "", errors.New("role not found in token")
}

// Injects user ID and role into request context (optional if using middleware)
func InjectUserContext(r *http.Request) *http.Request {
	tokenString, err := ExtractToken(r)
	if err != nil {
		return r
	}

	claims, err := ParseJWT(tokenString)
	if err != nil {
		return r
	}

	ctx := r.Context()
	if userID, ok := claims["user_id"].(string); ok {
		ctx = context.WithValue(ctx, "userId", userID)
	}	
	if role, ok := claims["role"].(string); ok {
		ctx = context.WithValue(ctx, "role", role)
	}

	return r.WithContext(ctx)
}
