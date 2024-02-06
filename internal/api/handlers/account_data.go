package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"notesAppBackend/internal/database"
	"strconv"
)

func AccountData(db *database.Database) gin.HandlerFunc {
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

		account, err := db.GetAccountData(accountId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, account)
	}
}
