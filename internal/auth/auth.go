package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Validate token secret
	if len(tokenSecret) == 0 {
		return "", errors.New("token secret cannot be empty")
	}

	// Signing key
	signingKey := []byte(tokenSecret)

	// Create claims
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Signing key
	signingKey := []byte(tokenSecret)

	// Validate token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	// Get subject claim from token
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	// Get user ID from subject claim
	userID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	// Get authorization header
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	// Check if authorization header is a bearer token
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("authorization header is not a bearer token")
	}

	// Check if bearer token is empty
	if len(authHeader[7:]) == 0 {
		return "", errors.New("bearer token is empty")
	}

	return authHeader[7:], nil
}

func MakeRefreshToken() (string, error) {
	// Generate random token
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}
