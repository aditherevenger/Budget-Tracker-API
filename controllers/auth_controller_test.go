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

type mockAuthService struct {
	RegisterFn      func(req *models.UserRegistrationRequest) (*models.UserResponse, error)
	LoginFn         func(req *models.UserLoginRequest) (string, *models.UserResponse, error)
	GetProfileFn    func(userID uint) (*models.UserResponse, error)
}

func (m *mockAuthService) Register(req *models.UserRegistrationRequest) (*models.UserResponse, error) {
	return m.RegisterFn(req)
}

func (m *mockAuthService) Login(req *models.UserLoginRequest) (string, *models.UserResponse, error) {
	return m.LoginFn(req)
}

func (m *mockAuthService) GetUserProfile(userID uint) (*models.UserResponse, error) {
	return m.GetProfileFn(userID)
}

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func performRequest(r http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestAuthController_Register_Success(t *testing.T) {
	mockSvc := &mockAuthService{
		RegisterFn: func(req *models.UserRegistrationRequest) (*models.UserResponse, error) {
			return &models.UserResponse{ID: 1, Email: req.Email, FirstName: req.FirstName, LastName: req.LastName}, nil
		},
	}
	ctrl := NewAuthController(mockSvc)
	r := setupGin()
	r.POST("/api/auth/register", ctrl.Register)

	payload := models.UserRegistrationRequest{Email: "jane@example.com", Password: "Pass1234", FirstName: "Jane", LastName: "Doe"}
	rec := performRequest(r, http.MethodPost, "/api/auth/register", payload, nil)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusCreated, rec.Code, rec.Body.String())
	}
}

func TestAuthController_Login_Success(t *testing.T) {
	mockSvc := &mockAuthService{
		LoginFn: func(req *models.UserLoginRequest) (string, *models.UserResponse, error) {
			return "token123", &models.UserResponse{ID: 1, Email: req.Email}, nil
		},
	}
	ctrl := NewAuthController(mockSvc)
	r := setupGin()
	r.POST("/api/auth/login", ctrl.Login)

	payload := models.UserLoginRequest{Email: "jane@example.com", Password: "Pass1234"}
	rec := performRequest(r, http.MethodPost, "/api/auth/login", payload, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestAuthController_Login_BadRequest(t *testing.T) {
	// No body => ShouldBindJSON fails -> 400
	mockSvc := &mockAuthService{}
	ctrl := NewAuthController(mockSvc)
	r := setupGin()
	r.POST("/api/auth/login", ctrl.Login)

	rec := performRequest(r, http.MethodPost, "/api/auth/login", nil, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func TestAuthController_GetProfile_Unauthorized(t *testing.T) {
	mockSvc := &mockAuthService{}
	ctrl := NewAuthController(mockSvc)
	r := setupGin()
	r.GET("/api/profile", ctrl.GetProfile)

	rec := performRequest(r, http.MethodGet, "/api/profile", nil, map[string]string{"Authorization": "Bearer sometoken"})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}

func TestAuthController_GetProfile_Success(t *testing.T) {
	mockSvc := &mockAuthService{
		GetProfileFn: func(userID uint) (*models.UserResponse, error) {
			return &models.UserResponse{ID: userID, Email: "jane@example.com"}, nil
		},
	}
	ctrl := NewAuthController(mockSvc)
	r := setupGin()
	r.GET("/api/profile", func(c *gin.Context) {
		c.Set("user_id", uint(42))
		ctrl.GetProfile(c)
	})

	rec := performRequest(r, http.MethodGet, "/api/profile", nil, map[string]string{"Authorization": "Bearer token"})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}
}
