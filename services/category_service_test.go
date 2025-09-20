package services

import (
	"errors"
	"testing"

	"github.com/aditherevenger/Budget-Tracker-API/models"
	"github.com/aditherevenger/Budget-Tracker-API/repository"
)

type mockCategoryRepo struct {
	CreateFn  func(category *models.Category) error
	ListFn    func(userID uint, filter *models.User) ([]models.Category, error)
	GetByIDFn func(id uint, userID uint) (*models.Category, error)
	UpdateFn  func(category *models.Category) error
	DeleteFn  func(id uint, userID uint) error
}

func (m *mockCategoryRepo) Create(category *models.Category) error                                  { return m.CreateFn(category) }
func (m *mockCategoryRepo) GetByUserID(userID uint, filter *models.User) ([]models.Category, error) { return m.ListFn(userID, filter) }
func (m *mockCategoryRepo) GetByID(id uint, userID uint) (*models.Category, error)                  { return m.GetByIDFn(id, userID) }
func (m *mockCategoryRepo) Update(category *models.Category) error                                  { return m.UpdateFn(category) }
func (m *mockCategoryRepo) Delete(id uint, userID uint) error                                       { return m.DeleteFn(id, userID) }

var _ repository.CategoryRepository = (*mockCategoryRepo)(nil)

func TestCategoryService_Create_DefaultColor(t *testing.T) {
	m := &mockCategoryRepo{ CreateFn: func(category *models.Category) error { category.ID = 1; return nil } }
	svc := NewCategoryService(m)
	cat, err := svc.CreateCategory(5, &models.CreateCategoryRequest{Name: "Food"})
	if err != nil { t.Fatalf("create: %v", err) }
	if cat.Color == "" { t.Fatalf("expected default color to be set, got empty") }
}

func TestCategoryService_GetCategories(t *testing.T) {
	m := &mockCategoryRepo{ ListFn: func(userID uint, filter *models.User) ([]models.Category, error) {
		return []models.Category{{ID: 1, UserID: userID, Name: "Food"}}, nil
	} }
	svc := NewCategoryService(m)
	cats, err := svc.GetCategories(7)
	if err != nil { t.Fatalf("get: %v", err) }
	if len(cats) != 1 || cats[0].Name != "Food" { t.Fatalf("unexpected: %+v", cats) }
}

func TestCategoryService_GetByID(t *testing.T) {
	m := &mockCategoryRepo{ GetByIDFn: func(id uint, userID uint) (*models.Category, error) { return &models.Category{ID: id, UserID: userID, Name: "Rent"}, nil } }
	svc := NewCategoryService(m)
	cat, err := svc.GetCategoryByID(2, 7)
	if err != nil { t.Fatalf("get: %v", err) }
	if cat.ID != 2 || cat.Name != "Rent" { t.Fatalf("unexpected: %+v", cat) }
}

func TestCategoryService_Update_PatchSemantics(t *testing.T) {
	m := &mockCategoryRepo{
		GetByIDFn: func(id uint, userID uint) (*models.Category, error) { return &models.Category{ID: id, UserID: userID, Name: "Old", Color: "#fff"}, nil },
		UpdateFn: func(category *models.Category) error { return nil },
	}
	svc := NewCategoryService(m)
	newName := "NewName"
	newColor := "#000"
	cat, err := svc.UpdateCategory(3, 7, &models.UpdateCategoryRequest{Name: &newName, Color: &newColor})
	if err != nil { t.Fatalf("update: %v", err) }
	if cat.Name != "NewName" || cat.Color != "#000" { t.Fatalf("unexpected: %+v", cat) }
}

func TestCategoryService_Delete(t *testing.T) {
	m := &mockCategoryRepo{ DeleteFn: func(id uint, userID uint) error { return nil } }
	svc := NewCategoryService(m)
	if err := svc.DeleteCategory(9, 7); err != nil { t.Fatalf("delete: %v", err) }
}

func TestCategoryService_Update_NotFound(t *testing.T) {
	m := &mockCategoryRepo{ GetByIDFn: func(id uint, userID uint) (*models.Category, error) { return nil, errors.New("not found") } }
	svc := NewCategoryService(m)
	if _, err := svc.UpdateCategory(3, 7, &models.UpdateCategoryRequest{}); err == nil {
		t.Fatalf("expected error when category not found")
	}
}
