package auth

import (
	"time"
	"testing"

	"github.com/google/uuid"
)

func TestJWT_RoundTrip(t *testing.T) {
	// Arrange
	secret := "super-secret-key-123"
	userID := uuid.New()
	duration := time.Hour

	// Act
	token, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate valid JWT: %v", err)
	}

	// Assert
	if parsedID != userID {
		t.Errorf("expected user ID %v, got %v", userID, parsedID)
	}
}

func TestJWT_ExpiredTokenRejected(t *testing.T) {
	// Arrange
	secret := "super-secret-key-123"
	userID := uuid.New()
	// Negative duration means the token is created in the past/expired immediately
	duration := -time.Minute 

	// Act
	token, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	_, err = ValidateJWT(token, secret)

	// Assert
	if err == nil {
		t.Error("expected error for expired token, but got nil")
	}
}

func TestJWT_WrongSecretRejected(t *testing.T) {
	// Arrange
	correctSecret := "correct-secret-key"
	wrongSecret := "wrong-secret-key"
	userID := uuid.New()
	duration := time.Hour

	// Act
	token, err := MakeJWT(userID, correctSecret, duration)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)

	// Assert
	if err == nil {
		t.Error("expected error for token validated with wrong secret, but got nil")
	}
}