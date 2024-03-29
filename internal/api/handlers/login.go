package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"notesAppBackend/internal/api"
	"notesAppBackend/internal/database"
	"notesAppBackend/internal/models"
	"notesAppBackend/internal/security"
)

// Login is responsible for authenticating an existing user.
func Login(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestJSON := struct {
			Nickname string `json:"nickname"`
			Password string `json:"password"`
		}{}

		if err := c.BindJSON(&requestJSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := db.GetUserByNickname(requestJSON.Nickname)
		if err != nil {
			api.HandleInternalServerError(c, err)
			return
		}

		if user == (models.User{}) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("User with name %s isn't found", requestJSON.Nickname),
			})
			return
		}

		hashedPassword, err := security.HashPassword(requestJSON.Password)
		if err != nil {
			api.HandleInternalServerError(c, err)
			return
		}

		if user.HashedPassword != hashedPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
			return
		}

		jwtToken, err := security.GenerateJWTToken(user)
		if err != nil {
			api.HandleInternalServerError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"jwt": jwtToken,
		})
	}
}
