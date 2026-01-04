package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// DefaultCost is the default bcrypt cost
const DefaultCost = bcrypt.DefaultCost

// HashPassword hashes a password using bcrypt with default cost
func HashPassword(password string) (string, error) {
	return HashPasswordWithCost(password, DefaultCost)
}

// HashPasswordWithCost hashes a password using bcrypt with custom cost
func HashPasswordWithCost(password string, cost int) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// CheckPasswordHash is an alias for VerifyPassword that returns bool
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
