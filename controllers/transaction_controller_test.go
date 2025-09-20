package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aditherevenger/Budget-Tracker-API/models"
	"github.com/gin-gonic/gin"
)

type mockTransactionService struct {
	CreateFn      func(userID uint, req *models.CreateTransactionRequest) (*models.Transaction, error)
	ListFn        func(userID uint, filter *models.TransactionFilter) ([]models.Transaction, error)
	GetByIDFn     func(id uint, userID uint) (*models.Transaction, error)
	UpdateFn      func(id uint, userID uint, req *models.UpdateTransactionRequest) (*models.Transaction, error)
	DeleteFn      func(id uint, userID uint) error
	SummaryFn     func(userID uint, startDate, endDate string) (map[string]interface{}, error)
}

func (m *mockTransactionService) CreateTransaction(userID uint, req *models.CreateTransactionRequest) (*models.Transaction, error) {
	return m.CreateFn(userID, req)
}
func (m *mockTransactionService) GetTransactions(userID uint, filter *models.TransactionFilter) ([]models.Transaction, error) {
	return m.ListFn(userID, filter)
}
func (m *mockTransactionService) GetTransactionByID(id uint, userID uint) (*models.Transaction, error) {
	return m.GetByIDFn(id, userID)
}
func (m *mockTransactionService) UpdateTransaction(id uint, userID uint, req *models.UpdateTransactionRequest) (*models.Transaction, error) {
	return m.UpdateFn(id, userID, req)
}
func (m *mockTransactionService) DeleteTransaction(id uint, userID uint) error { return m.DeleteFn(id, userID) }
func (m *mockTransactionService) GetSummary(userID uint, startDate, endDate string) (map[string]interface{}, error) {
	return m.SummaryFn(userID, startDate, endDate)
}

func setupGinTxn() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func performRequestTxn(r http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil { _ = json.NewEncoder(&buf).Encode(body) }
	req := httptest.NewRequest(method, path, &buf)
	for k, v := range headers { req.Header.Set(k, v) }
	if body != nil { req.Header.Set("Content-Type", "application/json") }
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestTransactionController_Create_Success(t *testing.T) {
	mockSvc := &mockTransactionService{ CreateFn: func(userID uint, req *models.CreateTransactionRequest) (*models.Transaction, error) {
		return &models.Transaction{ID: 1, UserID: userID, CategoryID: req.CategoryID, Amount: req.Amount, Type: req.Type, Description: req.Description, Date: req.Date}, nil
	}}
	ctrl := NewTransactionController(mockSvc)
	r := setupGinTxn()
	r.POST("/api/transactions/", func(c *gin.Context) { c.Set("user_id", uint(5)); ctrl.CreateTransaction(c) })

	payload := models.CreateTransactionRequest{CategoryID: 2, Amount: 10.5, Type: models.Expense, Description: "Coffee", Date: time.Now().UTC()}
	rec := performRequestTxn(r, http.MethodPost, "/api/transactions/", payload, nil)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d, body=%s", http.StatusCreated, rec.Code, rec.Body.String())
	}
}

func TestTransactionController_List_Success(t *testing.T) {
	mockSvc := &mockTransactionService{ ListFn: func(userID uint, filter *models.TransactionFilter) ([]models.Transaction, error) {
		return []models.Transaction{{ID: 1, UserID: userID, CategoryID: 2, Amount: 10.5, Type: models.Expense, Date: time.Now().UTC()}}, nil
	}}
	ctrl := NewTransactionController(mockSvc)
	r := setupGinTxn()
	r.GET("/api/transactions/", func(c *gin.Context) { c.Set("user_id", uint(5)); ctrl.GetTransactions(c) })

	rec := performRequestTxn(r, http.MethodGet, "/api/transactions/", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestTransactionController_GetByID_Success(t *testing.T) {
	mockSvc := &mockTransactionService{ GetByIDFn: func(id uint, userID uint) (*models.Transaction, error) {
		return &models.Transaction{ID: id, UserID: userID, CategoryID: 2, Amount: 10.5, Type: models.Expense, Date: time.Now().UTC()}, nil
	}}
	ctrl := NewTransactionController(mockSvc)
	r := setupGinTxn()
	r.GET("/api/transactions/:id", func(c *gin.Context) { c.Set("user_id", uint(5)); ctrl.GetTransaction(c) })

	rec := performRequestTxn(r, http.MethodGet, "/api/transactions/1", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestTransactionController_Update_Success(t *testing.T) {
	mockSvc := &mockTransactionService{ UpdateFn: func(id uint, userID uint, req *models.UpdateTransactionRequest) (*models.Transaction, error) {
		var amount float64 = 20
		if req.Amount != nil { amount = *req.Amount }
		return &models.Transaction{ID: id, UserID: userID, CategoryID: 2, Amount: amount, Type: models.Expense, Date: time.Now().UTC()}, nil
	}}
	ctrl := NewTransactionController(mockSvc)
	r := setupGinTxn()
	r.PUT("/api/transactions/:id", func(c *gin.Context) { c.Set("user_id", uint(5)); ctrl.UpdateTransaction(c) })

	amt := 22.75
	payload := models.UpdateTransactionRequest{Amount: &amt}
	rec := performRequestTxn(r, http.MethodPut, "/api/transactions/1", payload, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestTransactionController_Delete_Success(t *testing.T) {
	mockSvc := &mockTransactionService{ DeleteFn: func(id uint, userID uint) error { return nil } }
	ctrl := NewTransactionController(mockSvc)
	r := setupGinTxn()
	r.DELETE("/api/transactions/:id", func(c *gin.Context) { c.Set("user_id", uint(5)); ctrl.DeleteTransaction(c) })

	rec := performRequestTxn(r, http.MethodDelete, "/api/transactions/1", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestTransactionController_Summary_Success(t *testing.T) {
	mockSvc := &mockTransactionService{ SummaryFn: func(userID uint, startDate, endDate string) (map[string]interface{}, error) {
		return map[string]interface{}{"total_income": 1000.0, "total_expense": 200.0, "net": 800.0}, nil
	}}
	ctrl := NewTransactionController(mockSvc)
	r := setupGinTxn()
	r.GET("/api/transactions/summary", func(c *gin.Context) { c.Set("user_id", uint(5)); ctrl.GetSummary(c) })

	rec := performRequestTxn(r, http.MethodGet, "/api/transactions/summary?start_date=2025-01-01T00:00:00Z&end_date=2025-01-31T23:59:59Z", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d got %d, body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}
}

func TestTransactionController_Unauthorized_When_No_User(t *testing.T) {
	mockSvc := &mockTransactionService{}
	ctrl := NewTransactionController(mockSvc)
	r := setupGinTxn()
	r.GET("/api/transactions/", ctrl.GetTransactions)

	rec := performRequestTxn(r, http.MethodGet, "/api/transactions/", nil, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d got %d, body=%s", http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}
