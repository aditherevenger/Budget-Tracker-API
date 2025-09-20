package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestCategoryJSON_ContainsExpectedKeys(t *testing.T) {
	c := Category{
		ID:          11,
		UserID:      5,
		Name:        "Food",
		Description: "Meals & groceries",
		Color:       "#FF0000",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	b, err := json.Marshal(c)
	if err != nil { t.Fatalf("marshal error: %v", err) }
	js := string(b)

	for _, key := range []string{"\"id\"","\"user_id\"","\"name\"","\"description\"","\"color\"","\"created_at\"","\"updated_at\""} {
		if !strings.Contains(js, key) {
			t.Fatalf("expected JSON to contain %s, got: %s", key, js)
		}
	}
}

func TestCreateCategoryRequest_JSON_Marshalling(t *testing.T) {
	req := CreateCategoryRequest{Name: "Rent", Description: "Monthly rent", Color: "#00AAFF"}
	b, err := json.Marshal(req)
	if err != nil { t.Fatalf("marshal error: %v", err) }
	js := string(b)
	for _, key := range []string{"\"name\"","\"description\"","\"color\""} {
		if !strings.Contains(js, key) {
			t.Fatalf("expected JSON to contain %s, got: %s", key, js)
		}
	}
}

func TestUpdateCategoryRequest_OmitsNilFields(t *testing.T) {
	name := "Essentials"
	req := UpdateCategoryRequest{Name: &name}
	b, err := json.Marshal(req)
	if err != nil { t.Fatalf("marshal error: %v", err) }
	js := string(b)
	if !strings.Contains(js, "\"name\"") {
		t.Fatalf("expected name to be present, got: %s", js)
	}
	if strings.Contains(js, "\"description\"") || strings.Contains(js, "\"color\"") {
		t.Fatalf("expected nil fields to be omitted, got: %s", js)
	}
}
