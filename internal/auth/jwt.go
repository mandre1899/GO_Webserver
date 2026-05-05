package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetBearerToken(headers http.Header) (string, error) {
	jwtStr, ok := headers["Authorization"]
	if ok != false {
		return "", fmt.Errorf("Jwt notfound")
	}
	cleaned := strings.TrimPrefix(jwtStr[0], "Bearer ")
	return cleaned, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	ct := time.Now().UTC()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy-acces",
		IssuedAt:  jwt.NewNumericDate(ct),
		ExpiresAt: jwt.NewNumericDate(ct.Add(expiresIn)),
		Subject: userID.String(),
	})
	jwtStr, err := tok.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Couldn't sign the JWT: %w", err)
	}
	return jwtStr, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	if claims.Subject == "" {
		return uuid.Nil, errors.New("missing subject claim")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id in token subject: %w", err)
	}

	return userID, nil
}

