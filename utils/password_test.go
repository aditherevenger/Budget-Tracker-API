package utils

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	plain := "StrongPass123!"
	hash, err := HashPassword(plain)
	if err != nil { t.Fatalf("hash error: %v", err) }
	if hash == plain { t.Fatalf("hash should not equal plain text") }

	if !CheckPassword(plain, hash) {
		t.Fatalf("expected password to verify")
	}
	if CheckPassword("wrong", hash) {
		t.Fatalf("expected wrong password to fail")
	}
}
