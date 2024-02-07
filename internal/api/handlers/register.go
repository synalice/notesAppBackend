package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"notesAppBackend/internal/api"
	"notesAppBackend/internal/database"
	"notesAppBackend/internal/models"
	"notesAppBackend/internal/security"
	"time"
)

// Register is responsible for registering a new user.
func Register(db *database.Database) gin.HandlerFunc {
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

		if user != (models.User{}) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Nickname is already registered",
			})
			return
		}

		hashedPassword, err := security.HashPassword(requestJSON.Password)
		if err != nil {
			api.HandleInternalServerError(c, err)
			return
		}

		newUser := models.User{
			HashedPassword: hashedPassword,
			Nickname:       requestJSON.Nickname,
			DateCreated:    time.Now(),
		}

		newUser.ID, err = db.CreateNewUser(newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		jwtToken, err := security.GenerateJWTToken(newUser)
		if err != nil {
			api.HandleInternalServerError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"jwt": jwtToken,
		})
	}
}
