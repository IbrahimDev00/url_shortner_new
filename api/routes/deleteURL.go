package routes

import (
	"net/http"

	"github.com/AnuragChaubey/URL-Shortner/api/database"
	"github.com/gin-gonic/gin"
)

func DeleteURL(c *gin.Context) {
	shortID := c.Param("shortID")

	err := database.Client.Del(database.Ctx, shortID).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete shortened link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shortened link deleted successfully"})
}
