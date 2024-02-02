package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
	"net/http"
	"notesAppBackend/internal/api"
	"notesAppBackend/internal/security/secrets"
)

func VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestJSON api.JWTRequest
		if err := c.BindJSON(&requestJSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tokenString := requestJSON.JWT

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what it has to be.
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("VerifyJWT: unexpected signing method: %v", token.Method.Alg())
			}

			return []byte(secrets.JWTSecret), nil
		})

		if !token.Valid {
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token is malformed",
				})
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid token signature",
				})
			case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token is either expired or not active yet",
				})
			default:
				api.HandleInternalServerError(c, err)
			}
			return
		}

		c.Status(http.StatusOK)
	}
}
