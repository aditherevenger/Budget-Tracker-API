package utils

import (
	"testing"
	"time"
)

func TestGenerateAndValidateToken_Success(t *testing.T) {
	token, err := GenerateToken(7, "user@example.com")
	if err != nil { t.Fatalf("generate: %v", err) }
	claims, err := ValidateToken(token)
	if err != nil { t.Fatalf("validate: %v", err) }
	if claims.UserID != 7 || claims.Email != "user@example.com" { t.Fatalf("unexpected claims: %+v", claims) }
	if claims.ExpiresAt.Time.Before(time.Now()) { t.Fatalf("token should not be expired") }
}

func TestValidateToken_Invalid(t *testing.T) {
	_, err := ValidateToken("invalid.token.here")
	if err == nil { t.Fatalf("expected error for invalid token") }
}
