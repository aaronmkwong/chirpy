package auth

import "github.com/alexedwards/argon2id"

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
