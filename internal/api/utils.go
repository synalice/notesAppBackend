package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func HandleInternalServerError(c *gin.Context, err error) {
	log.Println(err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
