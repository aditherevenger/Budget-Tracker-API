package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestUserJSON_OmitsPasswordAndDeletedAt(t *testing.T) {
	u := User{
		ID:        10,
		Email:     "jane@example.com",
		Password:  "super-secret",
		FirstName: "Jane",
		LastName:  "Doe",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	b, err := json.Marshal(u)
	if err != nil { t.Fatalf("marshal error: %v", err) }
	js := string(b)

	if strings.Contains(js, "password") {
		t.Fatalf("expected password to be omitted from JSON, got: %s", js)
	}
	if !strings.Contains(js, "\"email\"") || !strings.Contains(js, "\"first_name\"") || !strings.Contains(js, "\"last_name\"") {
		t.Fatalf("expected key user fields to be present, got: %s", js)
	}
}

func TestUserResponseJSON_HasExpectedKeys(t *testing.T) {
	resp := UserResponse{ID: 1, Email: "jane@example.com", FirstName: "Jane", LastName: "Doe", CreatedAt: time.Now().UTC()}
	b, err := json.Marshal(resp)
	if err != nil { t.Fatalf("marshal error: %v", err) }
	js := string(b)
	for _, key := range []string{"\"id\"","\"email\"","\"first_name\"","\"last_name\"","\"created_at\""} {
		if !strings.Contains(js, key) {
			t.Fatalf("expected JSON to contain %s, got: %s", key, js)
		}
	}
}
