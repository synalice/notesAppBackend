package handlers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/golang-jwt/jwt/v5"
	"net/http"
	"notesAppBackend/internal/security"
)

func VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestJSON := struct {
			JWT string `json:"jwt"`
		}{}

		if err := c.BindJSON(&requestJSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := security.VerifyJWT(requestJSON.JWT)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}
