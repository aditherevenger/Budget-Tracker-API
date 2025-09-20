package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aditherevenger/Budget-Tracker-API/models"
	"github.com/gin-gonic/gin"
)

type mockCategoryService struct {
	CreateFn       func(userID uint, req *models.CreateCategoryRequest) (*models.Category, error)
	ListFn         func(userID uint) ([]models.Category, error)
	GetByIDFn      func(id uint, userID uint) (*models.Category, error)
	UpdateFn       func(id uint, userID uint, req *models.UpdateCategoryRequest) (*models.Category, error)
	DeleteFn       func(id uint, userID uint) error
}

func (m *mockCategoryService) CreateCategory(userID uint, req *models.CreateCategoryRequest) (*models.Category, error) {
	return m.CreateFn(userID, req)
}
func (m *mockCategoryService) GetCategories(userID uint) ([]models.Category, error) { return m.ListFn(userID) }
func (m *mockCategoryService) GetCategoryByID(id uint, userID uint) (*models.Category, error) {
	return m.GetByIDFn(id, userID)
}
func (m *mockCategoryService) UpdateCategory(id uint, userID uint, req *models.UpdateCategoryRequest) (*models.Category, error) {
	return m.UpdateFn(id, userID, req)
}
func (m *mockCategoryService) DeleteCategory(id uint, userID uint) error { return m.DeleteFn(id, userID) }

func setupGinCategory() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func performRequestCategory(r http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	for k, v := range headers { req.Header.Set(k, v) }
	if body != nil { req.Header.Set("Content-Type", "application/json") }
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestCategoryController_Create_Success(t *testing.T) {
	mockSvc := &mockCategoryService{ CreateFn: func(userID uint, req *models.CreateCategoryRequest) (*models.Category, error) {
		return &models.Category{ID: 1, UserID: userID, Name: req.Name}, nil
	}}
	ctrl := NewCategoryController(mockSvc)
	r := setupGinCategory()
	r.POST("/api/categories", func(c *gin.Context) { c.Set("user_id", uint(7)); ctrl.CreateCategories(c) })

	payload := models.CreateCategoryRequest{Name: "Food", Description: "Meals", Color: "#FF0000"}
	rec := performRequestCategory(r, http.MethodPost, "/api/categories", payload, nil)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d, body=%s", http.StatusCreated, rec.Code, rec.Body.String())
	}
}

func TestCategoryController_List_Success(t *testing.T) {
	mockSvc := &mockCategoryService{ ListFn: func(userID uint) ([]models.Category, error) {
		return []models.Category{{ID: 1, UserID: userID, Name: "Food"}}, nil
	}}
	ctrl := NewCategoryController(mockSvc)
	r := setupGinCategory()
	r.GET("/api/categories", func(c *gin.Context) { c.Set("user_id", uint(7)); ctrl.GetCategories(c) })

	rec := performRequestCategory(r, http.MethodGet, "/api/categories", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestCategoryController_Update_Success(t *testing.T) {
	mockSvc := &mockCategoryService{ UpdateFn: func(id uint, userID uint, req *models.UpdateCategoryRequest) (*models.Category, error) {
		name := "Updated"; if req.Name != nil { name = *req.Name }
		return &models.Category{ID: id, UserID: userID, Name: name}, nil
	}}
	ctrl := NewCategoryController(mockSvc)
	r := setupGinCategory()
	r.PUT("/api/categories/:id", func(c *gin.Context) { c.Set("user_id", uint(7)); ctrl.UpdateCategory(c) })

	newName := "Food & Dining"
	payload := models.UpdateCategoryRequest{Name: &newName}
	rec := performRequestCategory(r, http.MethodPut, "/api/categories/1", payload, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestCategoryController_Delete_Success(t *testing.T) {
	mockSvc := &mockCategoryService{ DeleteFn: func(id uint, userID uint) error { return nil } }
	ctrl := NewCategoryController(mockSvc)
	r := setupGinCategory()
	r.DELETE("/api/categories/:id", func(c *gin.Context) { c.Set("user_id", uint(7)); ctrl.DeleteCategory(c) })

	rec := performRequestCategory(r, http.MethodDelete, "/api/categories/1", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestCategoryController_Unauthorized_When_No_User(t *testing.T) {
	mockSvc := &mockCategoryService{}
	ctrl := NewCategoryController(mockSvc)
	r := setupGinCategory()
	r.GET("/api/categories", ctrl.GetCategories)

	rec := performRequestCategory(r, http.MethodGet, "/api/categories", nil, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d got %d, body=%s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}
