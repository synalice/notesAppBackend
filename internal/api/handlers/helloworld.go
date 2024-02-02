package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HelloWorld
// @Summary Ping an API.
// @Description This rout always responds with HTTP 200 and can be used
//
//	to verify that everything works.
//
// @Tags debug
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /helloworld [get]
func HelloWorld(g *gin.Context) {
	g.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
