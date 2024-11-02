package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{"ValidToken", uuid.New(), "mysecret", time.Hour, false},
		{"EmptySecret", uuid.New(), "", time.Hour, true},
		{"ZeroExpiration", uuid.New(), "mysecret", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(token) == 0 {
				t.Errorf("expected token to be non-empty")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	validUserID := uuid.New()
	validTokenSecret := "mysecret"
	validToken, _ := MakeJWT(validUserID, validTokenSecret, time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantErr     bool
		wantUserID  uuid.UUID
	}{
		{"ValidToken", validToken, validTokenSecret, false, validUserID},
		{"InvalidToken", "invalidtoken", validTokenSecret, true, uuid.Nil},
		{"WrongSecret", validToken, "wrongsecret", true, uuid.Nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
