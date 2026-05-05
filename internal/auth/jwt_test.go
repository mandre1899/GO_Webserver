package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	if token == "" {
		t.Fatal("MakeJWT returned empty token")
	}
	if len(token) == 0 {
		t.Error("Expected non-empty token")
	}
}

func TestMakeJWT_DifferentExpirations(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"

	testCases := []struct {
		name      string
		expiresIn time.Duration
	}{
		{"One hour", time.Hour},
		{"One day", 24 * time.Hour},
		{"One minute", time.Minute},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := MakeJWT(userID, tokenSecret, tc.expiresIn)
			if err != nil {
				t.Fatalf("MakeJWT failed: %v", err)
			}
			if token == "" {
				t.Error("Expected non-empty token")
			}
		})
	}
}

func TestValidateJWT_ValidToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	validatedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("Expected userID %s, got %s", userID, validatedUserID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	expiresIn := -time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Error("Expected error for expired token, but got nil")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	tokenSecret := "test-secret"

	testCases := []struct {
		name  string
		token string
	}{
		{"Empty string", ""},
		{"Invalid JWT", "invalid.token.here"},
		{"Gibberish", "xyz123nonsense"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidateJWT(tc.token, tokenSecret)
			if err == nil {
				t.Error("Expected error for invalid token")
			}
		})
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "original-secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	wrongSecret := "wrong-secret"
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("Expected error when validating with wrong secret")
	}
}

func TestMakeJWT_DifferentUsers(t *testing.T) {
	tokenSecret := "test-secret"
	expiresIn := time.Hour

	userID1 := uuid.New()
	userID2 := uuid.New()

	token1, err := MakeJWT(userID1, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed for user1: %v", err)
	}

	token2, err := MakeJWT(userID2, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed for user2: %v", err)
	}

	if token1 == token2 {
		t.Error("Expected different tokens for different users")
	}

	validatedID1, err := ValidateJWT(token1, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed for token1: %v", err)
	}
	if validatedID1 != userID1 {
		t.Errorf("Token1: Expected userID %s, got %s", userID1, validatedID1)
	}

	validatedID2, err := ValidateJWT(token2, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed for token2: %v", err)
	}
	if validatedID2 != userID2 {
		t.Errorf("Token2: Expected userID %s, got %s", userID2, validatedID2)
	}
}

func TestValidateJWT_InvalidSecretFormat(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	emptySecret := ""
	_, err = ValidateJWT(token, emptySecret)
	if err == nil {
		t.Error("Expected error when validating with empty secret")
	}
}

