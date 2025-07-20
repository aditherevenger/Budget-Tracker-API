package models

import (
	"encoding/json"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestUserStruct(t *testing.T) {
	t.Run("User struct creation and field assignment", func(t *testing.T) {
		now := time.Now()
		user := User{
			ID:        1,
			Email:     "test@example.com",
			Password:  "hashedpassword123",
			FirstName: "John",
			LastName:  "Doe",
			CreateAt:  now,
			UpdateAt:  now,
		}

		// Test field assignments
		if user.ID != 1 {
			t.Errorf("Expected ID to be 1, got %d", user.ID)
		}
		if user.Email != "test@example.com" {
			t.Errorf("Expected Email to be 'test@example.com', got %s", user.Email)
		}
		if user.Password != "hashedpassword123" {
			t.Errorf("Expected Password to be 'hashedpassword123', got %s", user.Password)
		}
		if user.FirstName != "John" {
			t.Errorf("Expected FirstName to be 'John', got %s", user.FirstName)
		}
		if user.LastName != "Doe" {
			t.Errorf("Expected LastName to be 'Doe', got %s", user.LastName)
		}
		if user.CreateAt != now {
			t.Errorf("Expected CreateAt to be %v, got %v", now, user.CreateAt)
		}
		if user.UpdateAt != now {
			t.Errorf("Expected UpdateAt to be %v, got %v", now, user.UpdateAt)
		}
	})

	t.Run("User struct with relationships", func(t *testing.T) {
		user := User{
			ID:           1,
			Email:        "test@example.com",
			Transactions: []Transaction{},
			categories:   []Category{},
		}

		// Test relationships initialization
		if user.Transactions == nil {
			t.Error("Expected Transactions to be initialized")
		}
		if user.categories == nil {
			t.Error("Expected categories to be initialized")
		}
		if len(user.Transactions) != 0 {
			t.Errorf("Expected empty Transactions slice, got length %d", len(user.Transactions))
		}
		if len(user.categories) != 0 {
			t.Errorf("Expected empty categories slice, got length %d", len(user.categories))
		}
	})

	t.Run("User struct with DeletedAt", func(t *testing.T) {
		now := time.Now()
		user := User{
			ID:       1,
			Email:    "test@example.com",
			DeleteAt: gorm.DeletedAt{Time: now, Valid: true},
		}

		if !user.DeleteAt.Valid {
			t.Error("Expected DeleteAt to be valid")
		}
		if user.DeleteAt.Time != now {
			t.Errorf("Expected DeleteAt time to be %v, got %v", now, user.DeleteAt.Time)
		}
	})
}

func TestUserJSON(t *testing.T) {
	t.Run("User JSON marshaling", func(t *testing.T) {
		now := time.Now()
		user := User{
			ID:        1,
			Email:     "test@example.com",
			Password:  "hashedpassword123",
			FirstName: "John",
			LastName:  "Doe",
			CreateAt:  now,
			UpdateAt:  now,
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("Failed to marshal User to JSON: %v", err)
		}

		// Verify JSON contains expected fields
		jsonStr := string(jsonData)
		expectedFields := []string{
			`"id":1`,
			`"email":"test@example.com"`,
			`"password":"hashedpassword123"`,
			`"first_name":"John"`,
			`"last_name":"Doe"`,
		}

		for _, field := range expectedFields {
			if !contains(jsonStr, field) {
				t.Errorf("Expected JSON to contain %s, got %s", field, jsonStr)
			}
		}
	})

	t.Run("User JSON unmarshaling", func(t *testing.T) {
		jsonStr := `{
			"id": 1,
			"email": "test@example.com",
			"password": "hashedpassword123",
			"first_name": "John",
			"last_name": "Doe"
		}`

		var user User
		err := json.Unmarshal([]byte(jsonStr), &user)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON to User: %v", err)
		}

		if user.ID != 1 {
			t.Errorf("Expected ID to be 1, got %d", user.ID)
		}
		if user.Email != "test@example.com" {
			t.Errorf("Expected Email to be 'test@example.com', got %s", user.Email)
		}
		if user.FirstName != "John" {
			t.Errorf("Expected FirstName to be 'John', got %s", user.FirstName)
		}
		if user.LastName != "Doe" {
			t.Errorf("Expected LastName to be 'Doe', got %s", user.LastName)
		}
	})
}

func TestUserRegistrationRequest(t *testing.T) {
	t.Run("UserRegistrationRequest struct creation", func(t *testing.T) {
		req := UserRegistrationRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		}

		if req.Email != "test@example.com" {
			t.Errorf("Expected Email to be 'test@example.com', got %s", req.Email)
		}
		if req.Password != "password123" {
			t.Errorf("Expected Password to be 'password123', got %s", req.Password)
		}
		if req.FirstName != "John" {
			t.Errorf("Expected FirstName to be 'John', got %s", req.FirstName)
		}
		if req.LastName != "Doe" {
			t.Errorf("Expected LastName to be 'Doe', got %s", req.LastName)
		}
	})

	t.Run("UserRegistrationRequest JSON marshaling", func(t *testing.T) {
		req := UserRegistrationRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("Failed to marshal UserRegistrationRequest to JSON: %v", err)
		}

		jsonStr := string(jsonData)
		expectedFields := []string{
			`"email":"test@example.com"`,
			`"password":"password123"`,
			`"first_name":"John"`,
			`"last_name":"Doe"`,
		}

		for _, field := range expectedFields {
			if !contains(jsonStr, field) {
				t.Errorf("Expected JSON to contain %s, got %s", field, jsonStr)
			}
		}
	})

	t.Run("UserRegistrationRequest JSON unmarshaling", func(t *testing.T) {
		jsonStr := `{
			"email": "test@example.com",
			"password": "password123",
			"first_name": "John",
			"last_name": "Doe"
		}`

		var req UserRegistrationRequest
		err := json.Unmarshal([]byte(jsonStr), &req)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON to UserRegistrationRequest: %v", err)
		}

		if req.Email != "test@example.com" {
			t.Errorf("Expected Email to be 'test@example.com', got %s", req.Email)
		}
		if req.Password != "password123" {
			t.Errorf("Expected Password to be 'password123', got %s", req.Password)
		}
		if req.FirstName != "John" {
			t.Errorf("Expected FirstName to be 'John', got %s", req.FirstName)
		}
		if req.LastName != "Doe" {
			t.Errorf("Expected LastName to be 'Doe', got %s", req.LastName)
		}
	})

	t.Run("UserRegistrationRequest with empty fields", func(t *testing.T) {
		req := UserRegistrationRequest{}

		if req.Email != "" {
			t.Errorf("Expected empty Email, got %s", req.Email)
		}
		if req.Password != "" {
			t.Errorf("Expected empty Password, got %s", req.Password)
		}
		if req.FirstName != "" {
			t.Errorf("Expected empty FirstName, got %s", req.FirstName)
		}
		if req.LastName != "" {
			t.Errorf("Expected empty LastName, got %s", req.LastName)
		}
	})
}

func TestUserLoginRequest(t *testing.T) {
	t.Run("UserLoginRequest struct creation", func(t *testing.T) {
		req := UserLoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		if req.Email != "test@example.com" {
			t.Errorf("Expected Email to be 'test@example.com', got %s", req.Email)
		}
		if req.Password != "password123" {
			t.Errorf("Expected Password to be 'password123', got %s", req.Password)
		}
	})

	t.Run("UserLoginRequest JSON marshaling", func(t *testing.T) {
		req := UserLoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("Failed to marshal UserLoginRequest to JSON: %v", err)
		}

		jsonStr := string(jsonData)
		expectedFields := []string{
			`"email":"test@example.com"`,
			`"password":"password123"`,
		}

		for _, field := range expectedFields {
			if !contains(jsonStr, field) {
				t.Errorf("Expected JSON to contain %s, got %s", field, jsonStr)
			}
		}
	})

	t.Run("UserLoginRequest JSON unmarshaling", func(t *testing.T) {
		jsonStr := `{
			"email": "test@example.com",
			"password": "password123"
		}`

		var req UserLoginRequest
		err := json.Unmarshal([]byte(jsonStr), &req)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON to UserLoginRequest: %v", err)
		}

		if req.Email != "test@example.com" {
			t.Errorf("Expected Email to be 'test@example.com', got %s", req.Email)
		}
		if req.Password != "password123" {
			t.Errorf("Expected Password to be 'password123', got %s", req.Password)
		}
	})

	t.Run("UserLoginRequest with empty fields", func(t *testing.T) {
		req := UserLoginRequest{}

		if req.Email != "" {
			t.Errorf("Expected empty Email, got %s", req.Email)
		}
		if req.Password != "" {
			t.Errorf("Expected empty Password, got %s", req.Password)
		}
	})
}

func TestUserResponse(t *testing.T) {
	t.Run("UserResponse struct creation", func(t *testing.T) {
		now := time.Now()
		resp := UserResponse{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			CreateAt:  now,
		}

		if resp.ID != 1 {
			t.Errorf("Expected ID to be 1, got %d", resp.ID)
		}
		if resp.Email != "test@example.com" {
			t.Errorf("Expected Email to be 'test@example.com', got %s", resp.Email)
		}
		if resp.FirstName != "John" {
			t.Errorf("Expected FirstName to be 'John', got %s", resp.FirstName)
		}
		if resp.LastName != "Doe" {
			t.Errorf("Expected LastName to be 'Doe', got %s", resp.LastName)
		}
		if resp.CreateAt != now {
			t.Errorf("Expected CreateAt to be %v, got %v", now, resp.CreateAt)
		}
	})

	t.Run("UserResponse JSON marshaling", func(t *testing.T) {
		now := time.Now()
		resp := UserResponse{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			CreateAt:  now,
		}

		jsonData, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("Failed to marshal UserResponse to JSON: %v", err)
		}

		jsonStr := string(jsonData)
		expectedFields := []string{
			`"id":1`,
			`"email":"test@example.com"`,
			`"first_name":"John"`,
			`"last_name":"Doe"`,
		}

		for _, field := range expectedFields {
			if !contains(jsonStr, field) {
				t.Errorf("Expected JSON to contain %s, got %s", field, jsonStr)
			}
		}
	})

	t.Run("UserResponse JSON unmarshaling", func(t *testing.T) {
		jsonStr := `{
			"id": 1,
			"email": "test@example.com",
			"first_name": "John",
			"last_name": "Doe"
		}`

		var resp UserResponse
		err := json.Unmarshal([]byte(jsonStr), &resp)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON to UserResponse: %v", err)
		}

		if resp.ID != 1 {
			t.Errorf("Expected ID to be 1, got %d", resp.ID)
		}
		if resp.Email != "test@example.com" {
			t.Errorf("Expected Email to be 'test@example.com', got %s", resp.Email)
		}
		if resp.FirstName != "John" {
			t.Errorf("Expected FirstName to be 'John', got %s", resp.FirstName)
		}
		if resp.LastName != "Doe" {
			t.Errorf("Expected LastName to be 'Doe', got %s", resp.LastName)
		}
	})

	t.Run("UserResponse with zero values", func(t *testing.T) {
		resp := UserResponse{}

		if resp.ID != 0 {
			t.Errorf("Expected ID to be 0, got %d", resp.ID)
		}
		if resp.Email != "" {
			t.Errorf("Expected empty Email, got %s", resp.Email)
		}
		if resp.FirstName != "" {
			t.Errorf("Expected empty FirstName, got %s", resp.FirstName)
		}
		if resp.LastName != "" {
			t.Errorf("Expected empty LastName, got %s", resp.LastName)
		}
	})
}

func TestStructComparisons(t *testing.T) {
	t.Run("User struct equality", func(t *testing.T) {
		now := time.Now()
		user1 := User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			CreateAt:  now,
		}
		user2 := User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			CreateAt:  now,
		}

		if user1.ID != user2.ID {
			t.Error("Expected users to have same ID")
		}
		if user1.Email != user2.Email {
			t.Error("Expected users to have same Email")
		}
		if user1.FirstName != user2.FirstName {
			t.Error("Expected users to have same FirstName")
		}
		if user1.LastName != user2.LastName {
			t.Error("Expected users to have same LastName")
		}
	})

	t.Run("UserRegistrationRequest struct equality", func(t *testing.T) {
		req1 := UserRegistrationRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		}
		req2 := UserRegistrationRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		}

		if req1 != req2 {
			t.Error("Expected UserRegistrationRequest structs to be equal")
		}
	})

	t.Run("UserLoginRequest struct equality", func(t *testing.T) {
		req1 := UserLoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		req2 := UserLoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		if req1 != req2 {
			t.Error("Expected UserLoginRequest structs to be equal")
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Mock types for testing relationships (since they're referenced but not defined)
type Transaction struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
}

type Category struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
}
