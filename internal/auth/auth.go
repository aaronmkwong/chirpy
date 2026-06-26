package auth

import (
	"time"
	"fmt"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/alexedwards/argon2id"
)

// HashPassword hashes a plaintext password using Argon2id.
func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(
	password,
	argon2id.DefaultParams,
	)
}

// CheckPasswordHash compares a plaintext password against a stored hash.
func CheckPasswordHash(password string, hash string,) (bool, error) {
	return argon2id.ComparePasswordAndHash(
	password,
	hash,
	)
}

// MakeJWT creates a signed JWT for the given user that expires after expiresIn.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// tokenSecret must be of type []byte for HMAC SHA-256
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateJWT verifies a signed JWT and returns the user ID stored in its subject claim.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{} // empty form: creates an empty claims struct and takes a pointer to it

	// Parse and validate the token string
	token, err := jwt.ParseWithClaims( // library fills in "form" and now claims.Subject, claims.ExpiresAt, etc. hold the token's data
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			// Validate that the signing method matches the expected HMAC method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(tokenSecret), nil
		},
	)

	if err != nil {
		return uuid.Nil, err
	}

	// Double-check token validity (jwt.ParseWithClaims handles expiration and signature checks)
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	// Extract the user ID from the Subject field
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	// Parse the stringified Subject back into a uuid.UUID
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID format in token: %w", err)
	}

	return userID, nil
}