package security

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"notesAppBackend/internal/models"
	"notesAppBackend/internal/security/secrets"
	"strconv"
	"time"
)

type JWTClaims struct {
	HeyThere string `json:"heyThere"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(user models.User) (string, error) {
	expirationDate := time.Now().Add(time.Hour * 24 * 2) // Expire token after 2 days
	id := uuid.New()

	claims := JWTClaims{
		HeyThere: "Get out of this JWT!",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "notesAppBackend",
			Subject:   strconv.Itoa(user.ID),
			ExpiresAt: jwt.NewNumericDate(expirationDate),
			ID:        id.String(),
		},
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	signedToken, err := unsignedToken.SignedString([]byte(secrets.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("GenerateJWTToken: %w", err)
	}

	return signedToken, nil
}
