package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"notesAppBackend/internal/api"
	"notesAppBackend/internal/database"
	"strconv"
)

// UserExists returns 200 is user exists and 404 is not.
func UserExists(db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountIdFromQuery, ok := c.GetQuery("id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No query is present"})
			return
		}

		accountId, err := strconv.Atoi(accountIdFromQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be an integer"})
			return
		}

		exists, err := db.IsIDPresent(accountId)
		if err != nil {
			api.HandleInternalServerError(c, err)
		}

		if exists {
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusNotFound)
		}
	}
}
