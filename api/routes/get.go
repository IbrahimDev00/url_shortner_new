package routes

import (
	"net/http"

	"github.com/AnuragChaubey/URL-Shortner/api/database"
	"github.com/gin-gonic/gin"
)

func GetByShortID(c *gin.Context) {
	shortID := c.Param("shortID")

	val, err := database.Client.Get(database.Ctx, shortID).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found for the given tagID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": val})
}
