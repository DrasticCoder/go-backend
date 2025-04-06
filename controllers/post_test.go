package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"go-backend/config"
	"go-backend/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestPost(authorID primitive.ObjectID, status string, scheduledAt *time.Time) models.Post {
	post := models.Post{
		ID:          primitive.NewObjectID(),
		AuthorID:    authorID,
		Title:       "Test Post",
		Content:     "Sample content",
		Type:        "trade",
		Visibility:  "public",
		Tags:        []string{"test"},
		Status:      status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ScheduledAt: scheduledAt,
	}

	if status == "published" && scheduledAt != nil {
		post.PublishedAt = scheduledAt
	}

	_, _ = config.PostCollection.InsertOne(context.TODO(), post)
	return post
}

func TestCreatePost(t *testing.T) {
	authorID := primitive.NewObjectID()

	payload := models.Post{
		Title:      "Draft Post Test",
		Content:    "Hello from test",
		Type:       "idea",
		Visibility: "public",
		Tags:       []string{"test"},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
	req = injectTestAuth(req, authorID.Hex(), "author")
	w := httptest.NewRecorder()

	CreatePost(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestListPosts(t *testing.T) {
	authorID := primitive.NewObjectID()
	now := time.Now().Add(-time.Hour)
	setupTestPost(authorID, "published", &now)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
	w := httptest.NewRecorder()

	ListPosts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestListMyPosts(t *testing.T) {
	authorID := primitive.NewObjectID()
	setupTestPost(authorID, "draft", nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/me", nil)
	req = injectTestAuth(req, authorID.Hex(), "author")
	w := httptest.NewRecorder()

	ListMyPosts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestUpdatePost(t *testing.T) {
	authorID := primitive.NewObjectID()
	post := setupTestPost(authorID, "draft", nil)

	update := models.Post{
		Title: "Updated Title from test",
	}
	body, _ := json.Marshal(update)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/posts/"+post.ID.Hex(), bytes.NewReader(body))
	req = muxWithParams(req, "id", post.ID.Hex())
	req = injectTestAuth(req, authorID.Hex(), "author")
	w := httptest.NewRecorder()

	UpdatePost(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

func TestDeletePost(t *testing.T) {
	authorID := primitive.NewObjectID()
	post := setupTestPost(authorID, "draft", nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/"+post.ID.Hex(), nil)
	req = muxWithParams(req, "id", post.ID.Hex())
	req = injectTestAuth(req, authorID.Hex(), "author")
	w := httptest.NewRecorder()

	DeletePost(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}


func injectTestAuth(req *http.Request, userID string, role string) *http.Request {
	ctx := context.WithValue(req.Context(), "userID", userID)
	ctx = context.WithValue(ctx, "role", role)
	return req.WithContext(ctx)
}


func muxWithParams(req *http.Request, key, value string) *http.Request {
	vars := map[string]string{key: value}
	return mux.SetURLVars(req, vars)
}
