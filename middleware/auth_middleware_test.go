package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aditherevenger/Budget-Tracker-API/utils"
	"github.com/gin-gonic/gin"
)

func setupRouterWithAuth() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		uid, _ := c.Get("user_id")
		email, _ := c.Get("user_email")
		c.JSON(http.StatusOK, gin.H{"user_id": uid, "email": email})
	})
	return r
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	r := setupRouterWithAuth()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestAuthMiddleware_MalformedBearer(t *testing.T) {
	r := setupRouterWithAuth()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token abc")
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	r := setupRouterWithAuth()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	r := setupRouterWithAuth()
	rec := httptest.NewRecorder()
	token, err := utils.GenerateToken(99, "valid@example.com")
	if err != nil { t.Fatalf("generate token: %v", err) }
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}
}
