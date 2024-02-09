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
	jwt.RegisteredClaims
}

func GenerateJWTToken(user models.User) (string, error) {
	expirationDate := time.Now().Add(time.Hour * 24 * 2) // Expire token after 2 days
	id := uuid.New()

	claims := JWTClaims{
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

// VerifyJWTValidity only checks if the JWT has correct structure, has not been tampered with
// and is not expired.
// TODO: Add check for if the token has expired or not.
func VerifyJWTValidity(token string) (JWTClaims, error) {
	t, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what it has to be.
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("VerifyJWTValidity: unexpected signing method: %v", token.Method.Alg())
		}

		return []byte(secrets.JWTSecret), nil
	})

	if err != nil || !t.Valid {
		return JWTClaims{}, fmt.Errorf("VerifyJWTValidity: %w", err)
	}

	return *t.Claims.(*JWTClaims), nil
}

func GetSubjectFromJWTClaims(claims JWTClaims) (int, error) {
	sub, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, fmt.Errorf("GetSubjectFromJWTClaims: %w", err)
	}
	return sub, nil
}
