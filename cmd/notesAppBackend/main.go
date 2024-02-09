package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"notesAppBackend/docs"
	"notesAppBackend/internal/api/handlers"
	"notesAppBackend/internal/database"
)

func main() {
	envs, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable", envs["PG_USER"], envs["PG_PASSWORD"], envs["PG_HOST"], envs["PG_PORT"])
	db, err := database.New(connStr)
	if err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	config := cors.DefaultConfig()
	config.AddAllowHeaders("Authorization")
	// TODO: Somehow remove "http://localhost:3000" in production. It's unsafe.
	config.AllowOrigins = []string{"http://localhost:3000", "https://syn-dev.ru/"}
	router.Use(cors.New(config))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/account-data", handlers.AccountData(db))
		v1.GET("/user-exists", handlers.UserExists(db))

		v1.POST("/register", handlers.Register(db))
		v1.POST("/login", handlers.Login(db))
		v1.POST("/verify-jwt", handlers.VerifyJWT())
	}

	post := v1.Group("/post")
	{
		post.POST("/new", handlers.NewPost(db))
		post.DELETE("/:id", handlers.DeleteNote(db))
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	_ = router.Run("0.0.0.0:8080")
}

// TODO: Do some refactoring.
// TODO: Handlers are leaking internal error o the client and don't have sufficient logging. Improve this!
// TODO: Rename every `post` to `note`.
