package services

import (
	"testing"

	"gorm.io/gorm"

	"github.com/aditherevenger/Budget-Tracker-API/models"
	"github.com/aditherevenger/Budget-Tracker-API/repository"
	"github.com/aditherevenger/Budget-Tracker-API/utils"
)

type mockUserRepo struct {
	CreateFn    func(user *models.User) error
	GetByEmailFn func(email string) (*models.User, error)
	GetByIDFn   func(id uint) (*models.User, error)
	UpdateFn    func(user *models.User) error
	DeleteFn    func(user *models.User) error
}

func (m *mockUserRepo) Create(user *models.User) error                     { return m.CreateFn(user) }
func (m *mockUserRepo) GetByEmail(email string) (*models.User, error)      { return m.GetByEmailFn(email) }
func (m *mockUserRepo) GetByID(id uint) (*models.User, error)              { return m.GetByIDFn(id) }
func (m *mockUserRepo) Update(user *models.User) error                     { return m.UpdateFn(user) }
func (m *mockUserRepo) Delete(user *models.User) error                     { return m.DeleteFn(user) }

var _ repository.UserRepository = (*mockUserRepo)(nil)

func TestAuthService_Register_Success(t *testing.T) {
	m := &mockUserRepo{
		GetByEmailFn: func(email string) (*models.User, error) { return nil, gorm.ErrRecordNotFound },
		CreateFn: func(user *models.User) error { user.ID = 1; return nil },
	}
	svc := NewAuthService(m)
	resp, err := svc.Register(&models.UserRegistrationRequest{Email: "jane@example.com", Password: "Pass1234", FirstName: "Jane", LastName: "Doe"})
	if err != nil { t.Fatalf("Register error: %v", err) }
	if resp.Email != "jane@example.com" || resp.ID == 0 { t.Fatalf("unexpected resp: %+v", resp) }
}

func TestAuthService_Register_Duplicate(t *testing.T) {
	m := &mockUserRepo{
		GetByEmailFn: func(email string) (*models.User, error) { return &models.User{ID: 99, Email: email}, nil },
	}
	svc := NewAuthService(m)
	if _, err := svc.Register(&models.UserRegistrationRequest{Email: "dup@example.com", Password: "x", FirstName: "A", LastName: "B"}); err == nil {
		t.Fatalf("expected error for duplicate email")
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	hashed, _ := utils.HashPassword("Pass1234")
	m := &mockUserRepo{
		GetByEmailFn: func(email string) (*models.User, error) { return &models.User{ID: 3, Email: email, Password: hashed, FirstName: "J", LastName: "D"}, nil },
	}
	svc := NewAuthService(m)
	token, user, err := svc.Login(&models.UserLoginRequest{Email: "jane@example.com", Password: "Pass1234"})
	if err != nil { t.Fatalf("login error: %v", err) }
	if token == "" { t.Fatalf("expected token, got empty") }
	if user.ID != 3 || user.Email != "jane@example.com" { t.Fatalf("unexpected user: %+v", user) }
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	hashed, _ := utils.HashPassword("CorrectPass")
	m := &mockUserRepo{ GetByEmailFn: func(email string) (*models.User, error) { return &models.User{ID: 3, Email: email, Password: hashed}, nil } }
	svc := NewAuthService(m)
	if _, _, err := svc.Login(&models.UserLoginRequest{Email: "jane@example.com", Password: "wrong"}); err == nil {
		t.Fatalf("expected invalid credentials error")
	}
}

func TestAuthService_GetUserProfile_Success(t *testing.T) {
	m := &mockUserRepo{ GetByIDFn: func(id uint) (*models.User, error) { return &models.User{ID: id, Email: "jane@example.com", FirstName: "Jane", LastName: "Doe"}, nil } }
	svc := NewAuthService(m)
	resp, err := svc.GetUserProfile(42)
	if err != nil { t.Fatalf("GetUserProfile error: %v", err) }
	if resp.ID != 42 || resp.Email != "jane@example.com" { t.Fatalf("unexpected resp: %+v", resp) }
}
