package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"notesAppBackend/docs"
	"notesAppBackend/internal/api/handlers"
	"notesAppBackend/internal/database"
)

func main() {
	// TODO: Get password from .env file instead!
	connStr := "postgres://postgres:pbnppl44@postgres:5432/postgres?sslmode=disable"
	db, err := database.New(connStr)
	if err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	config := cors.DefaultConfig()
	config.AddAllowHeaders("Authorization")
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/account-data", handlers.AccountData(db))
		v1.GET("/user-exists", handlers.UserExists(db))
		v1.POST("/register", handlers.Register(db))
		v1.POST("/login", handlers.Login(db))
		v1.POST("/verify-jwt", handlers.VerifyJWT())
		v1.POST("/new-post", handlers.NewPost(db))
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	_ = router.Run("0.0.0.0:8080")
}

// TODO: Do some refactoring.
