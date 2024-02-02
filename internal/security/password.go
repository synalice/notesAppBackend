package security

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"log"
	"notesAppBackend/internal/security/secrets"
)

// HashPassword uses scrypt algorithm to produce a hash of the password
func HashPassword(password string) (string, error) {
	// Read more about these magic numbers in the package's documentation.
	hashedPass, err := scrypt.Key([]byte(password), []byte(secrets.PasswordSalt), 32768, 8, 1, 32)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("could not hash the password")
	}

	return base64.StdEncoding.EncodeToString(hashedPass), nil
}
