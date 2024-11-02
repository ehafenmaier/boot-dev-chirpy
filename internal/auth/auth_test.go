package auth

import (
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"ValidPassword", "mysecretpassword", false},
		{"EmptyPassword", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(hash) == 0 {
				t.Errorf("expected hash to be non-empty")
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{"ValidPassword", "mysecretpassword", func() string { h, _ := HashPassword("mysecretpassword"); return h }(), false},
		{"WrongPassword", "wrongpassword", func() string { h, _ := HashPassword("mysecretpassword"); return h }(), true},
		{"EmptyPassword", "", func() string { h, _ := HashPassword("mysecretpassword"); return h }(), true},
		{"EmptyHash", "mysecretpassword", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name:    "ValidBearerToken",
			headers: http.Header{"Authorization": []string{"Bearer validtoken"}},
			want:    "validtoken",
			wantErr: false,
		},
		{
			name:    "MissingAuthorizationHeader",
			headers: http.Header{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "InvalidBearerToken",
			headers: http.Header{"Authorization": []string{"InvalidToken"}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "EmptyBearerToken",
			headers: http.Header{"Authorization": []string{"Bearer "}},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
