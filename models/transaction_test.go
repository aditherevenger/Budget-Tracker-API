package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestTransactionJSON_ContainsExpectedKeys(t *testing.T) {
	dt := time.Date(2025, 9, 18, 0, 0, 0, 0, time.UTC)
	tx := Transaction{
		ID:         21,
		UserID:     7,
		CategoryID: 3,
		Amount:     125.50,
		Type:       Expense,
		Description:"Groceries",
		Date:       dt,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	b, err := json.Marshal(tx)
	if err != nil { t.Fatalf("marshal error: %v", err) }
	js := string(b)

	for _, key := range []string{"\"id\"","\"user_id\"","\"category_id\"","\"amount\"","\"type\"","\"description\"","\"date\"","\"created_at\"","\"updated_at\""} {
		if !strings.Contains(js, key) {
			t.Fatalf("expected JSON to contain %s, got: %s", key, js)
		}
	}
}

func TestCreateTransactionRequest_JSON_Marshalling(t *testing.T) {
	dt := time.Date(2025, 9, 18, 0, 0, 0, 0, time.UTC)
	req := CreateTransactionRequest{CategoryID: 3, Amount: 10.0, Type: Income, Description: "Salary", Date: dt}
	b, err := json.Marshal(req)
	if err != nil { t.Fatalf("marshal error: %v", err) }
	js := string(b)
	for _, key := range []string{"\"category_id\"","\"amount\"","\"type\"","\"description\"","\"date\""} {
		if !strings.Contains(js, key) {
			t.Fatalf("expected JSON to contain %s, got: %s", key, js)
		}
	}
}

func TestUpdateTransactionRequest_OmitsNilFields(t *testing.T) {
	amt := 99.99
	typ := Expense
	dt := time.Date(2025, 9, 20, 0, 0, 0, 0, time.UTC)
	req := UpdateTransactionRequest{Amount: &amt, Type: &typ, Date: &dt}
	b, err := json.Marshal(req)
	if err != nil { t.Fatalf("marshal error: %v", err) }
	js := string(b)
	for _, key := range []string{"\"amount\"","\"type\"","\"date\""} {
		if !strings.Contains(js, key) {
			t.Fatalf("expected JSON to contain %s, got: %s", key, js)
		}
	}
	if strings.Contains(js, "\"category_id\"") || strings.Contains(js, "\"description\"") {
		t.Fatalf("expected nil fields to be omitted, got: %s", js)
	}
}
