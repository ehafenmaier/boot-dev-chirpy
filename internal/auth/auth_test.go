package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(hash) == 0 {
		t.Fatalf("expected hash to be non-empty")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mysecretpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wrongPassword := "wrongpassword"
	err = CheckPasswordHash(wrongPassword, hash)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
