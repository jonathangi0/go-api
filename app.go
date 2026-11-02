package main

import (
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"	
	"github.com/gothinkster/golang-gin-realworld-example-app/common"
	"github.com/gothinkster/golang-gin-realworld-example-app/intranet/branch/nvo"
	"github.com/gothinkster/golang-gin-realworld-example-app/users"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	users.AutoMigrate()
	/* db.AutoMigrate(&articles.ArticleModel{})
	db.AutoMigrate(&articles.TagModel{})
	db.AutoMigrate(&articles.FavoriteModel{})
	db.AutoMigrate(&articles.ArticleUserModel{})
	db.AutoMigrate(&articles.CommentModel{}) */
}

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal("failed to load .env:", err)
	}

	db, err := common.Init()
	if err != nil {
		log.Fatal("failed to initialize database:", err)
	}
	//Migrate(db)
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("failed to get sql.DB:", err)
	} else {
		defer sqlDB.Close()
	}

	r := gin.Default()

	// Disable automatic redirect for trailing slashes
	// This prevents POST body from being lost during redirects
	r.RedirectTrailingSlash = false

	v1 := r.Group("/api")
	//users.UsersRegister(v1.Group("/users"))
	v1.Use(users.AuthMiddleware(true))
	//articles.ArticlesAnonymousRegister(v1.Group("/articles"))
	//articles.TagsAnonymousRegister(v1.Group("/tags"))
	//users.ProfileRetrieveRegister(v1.Group("/profiles"))

	//v1.Use(users.AuthMiddleware(true))
	//users.UserRegister(v1.Group("/user"))
	//users.ProfileRegister(v1.Group("/profiles"))

	//articles.ArticlesRegister(v1.Group("/articles"))
	nvo.SurgicalNotesRegister(v1.Group("/surgical-notes"))

	// Intranet 2026 Migrate to Laravel a Go

	testAuth := r.Group("/api/ping")

	testAuth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
