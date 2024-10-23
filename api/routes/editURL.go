package routes

import (
	"net/http"
	"time"

	"github.com/AnuragChaubey/URL-Shortner/api/database"
	"github.com/AnuragChaubey/URL-Shortner/api/models"
	"github.com/gin-gonic/gin"
)

func EditURL(c *gin.Context) {
	shortID := c.Param("shortID")
	var body models.Request

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot parse JSON"})
		return
	}

	// Check if the shortID exists in the database
	val, err := database.Client.Get(database.Ctx, shortID).Result()
	if err != nil || val == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShortID does not exist"})
		return
	}

	// Update the URL associated with the shortID
	err = database.Client.Set(database.Ctx, shortID, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update shortened link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shortened link updated successfully"})
}
