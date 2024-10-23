package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AnuragChaubey/URL-Shortner/api/database"
	"github.com/AnuragChaubey/URL-Shortner/api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	database.InitializeClient()

	router := gin.Default()

	setupRoutes(router)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(router.Run(":" + port))
}

func setupRoutes(router *gin.Engine) {
	router.POST("/api/v1", routes.ShortenURL)
	router.PUT("/api/v1/:shortID", routes.EditURL)
	router.DELETE("/api/v1/:shortID", routes.DeleteURL)

	// Tag operations
	router.POST("/api/v1/addTag", routes.AddTag)

	// New GET endpoint for fetching data based on tag ID
	router.GET("/api/v1/:shortID", routes.GetByShortID)
}
