package jwt

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Set JWT_SECRET for testing
	os.Setenv("JWT_SECRET", "test-secret-key")
	code := m.Run()
	os.Exit(code)
}

func TestGenerateToken_Success(t *testing.T) {
	user := UserPayload{
		ID:       "user-123",
		Email:    "test@example.com",
		Username: "testuser",
	}

	token, err := GenerateToken(user)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateTokenWithExpiry_Success(t *testing.T) {
	user := UserPayload{
		ID:       "user-123",
		Email:    "test@example.com",
		Username: "testuser",
	}
	expiry := 1 * time.Hour

	token, err := GenerateTokenWithExpiry(user, expiry)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken_Success(t *testing.T) {
	secret := "test-secret-key"
	user := UserPayload{
		ID:       "user-123",
		Email:    "test@example.com",
		Username: "testuser",
	}

	token, err := GenerateToken(user)
	require.NoError(t, err)

	claims, err := ValidateToken(secret, token)

	require.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Username, claims.Username)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	secret := "test-secret-key"
	invalidToken := "invalid.token.here"

	claims, err := ValidateToken(secret, invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	wrongSecret := "wrong-secret"
	user := UserPayload{
		ID:       "user-123",
		Email:    "test@example.com",
		Username: "testuser",
	}

	token, err := GenerateToken(user)
	require.NoError(t, err)

	claims, err := ValidateToken(wrongSecret, token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	secret := "test-secret-key"
	user := UserPayload{
		ID:       "user-123",
		Email:    "test@example.com",
		Username: "testuser",
	}

	// Generate token with very short expiry
	token, err := GenerateTokenWithExpiry(user, -1*time.Hour) // Already expired
	require.NoError(t, err)

	claims, err := ValidateToken(secret, token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}
