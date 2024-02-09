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

func DeleteNote(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		noteID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect note ID"})
			return
		}

		jwtToken := strings.Split(c.GetHeader("Authorization"), " ")[1]
		jwtClaims, err := security.VerifyJWTValidity(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := security.GetSubjectFromJWTClaims(jwtClaims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		note, err := db.FindNoteByID(noteID)
		if err != nil {
			api.HandleInternalServerError(c, err)
			return
		}

		if userID != note.AuthorID {
			c.Status(http.StatusUnauthorized)
			return
		}

		err = db.DeleteNoteByID(noteID)
		if err != nil {
			api.HandleInternalServerError(c, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
