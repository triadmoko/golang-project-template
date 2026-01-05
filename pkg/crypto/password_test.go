package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword_Success(t *testing.T) {
	password := "testPassword123"

	hashedPassword, err := HashPassword(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
}

func TestHashPasswordWithCost_Success(t *testing.T) {
	password := "testPassword123"
	cost := bcrypt.MinCost

	hashedPassword, err := HashPasswordWithCost(password, cost)

	require.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)
}

func TestHashPasswordWithCost_InvalidCost(t *testing.T) {
	password := "testPassword123"
	cost := 100 // Invalid cost (too high)

	_, err := HashPasswordWithCost(password, cost)

	assert.Error(t, err)
}

func TestVerifyPassword_Success(t *testing.T) {
	password := "testPassword123"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	err = VerifyPassword(hashedPassword, password)

	assert.NoError(t, err)
}

func TestVerifyPassword_Invalid(t *testing.T) {
	password := "testPassword123"
	wrongPassword := "wrongPassword"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	err = VerifyPassword(hashedPassword, wrongPassword)

	assert.Error(t, err)
}

func TestCheckPasswordHash_True(t *testing.T) {
	password := "testPassword123"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	result := CheckPasswordHash(password, hashedPassword)

	assert.True(t, result)
}

func TestCheckPasswordHash_False(t *testing.T) {
	password := "testPassword123"
	wrongPassword := "wrongPassword"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	result := CheckPasswordHash(wrongPassword, hashedPassword)

	assert.False(t, result)
}
