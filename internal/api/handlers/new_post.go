package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"notesAppBackend/internal/api"
	"notesAppBackend/internal/database"
	"notesAppBackend/internal/security"
	"strconv"
	"strings"
)

func NewPost(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestJSON := struct {
			Title    string `json:"title"`
			Contents string `json:"contents"`
		}{}

		if err := c.BindJSON(&requestJSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't parse request body"})
			return
		}

		jwtToken := strings.Split(c.GetHeader("Authorization"), " ")[1]
		jwtClaims, err := security.VerifyJWTValidity(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		authorID, err := strconv.Atoi(jwtClaims.Subject)
		if err != nil {
			api.HandleInternalServerError(c, err)
			return
		}

		err = db.CreateNewPost(authorID, requestJSON.Title, requestJSON.Contents)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	}
}
