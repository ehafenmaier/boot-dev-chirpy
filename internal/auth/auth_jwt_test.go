package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(token) == 0 {
		t.Fatalf("expected token to be non-empty")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	validatedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if validatedUserID != userID {
		t.Fatalf("expected userID %v, got %v", userID, validatedUserID)
	}

	// Test with invalid token
	_, err = ValidateJWT("invalidtoken", tokenSecret)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
