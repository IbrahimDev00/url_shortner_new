package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AnuragChaubey/URL-Shortner/api/database"
	"github.com/AnuragChaubey/URL-Shortner/api/helpers"
	"github.com/AnuragChaubey/URL-Shortner/api/models"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func ShortenURL(c *gin.Context) {
	var body models.Request

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot Parse JSON"})
		return
	}

	val, err := database.Client.Get(database.Ctx, c.ClientIP()).Result()
	if err == redis.Nil {
		_ = database.Client.Set(database.Ctx, c.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val, _ = database.Client.Get(database.Ctx, c.ClientIP()).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := database.Client.TTL(database.Ctx, c.ClientIP()).Result()
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":            "rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
			return
		}
	}

	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	if !helpers.IsDifferentDomain(body.URL) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "You Can't Hack this System (:",
		})
		return
	}

	body.URL = helpers.EnsureHTTPPrefix(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	val, _ = database.Client.Get(database.Ctx, id).Result()
	if val != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "URL Custom Short is already in use",
		})
		return
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = database.Client.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to connect to server",
		})
		return
	}

	resp := models.Response{
		Expiry:          body.Expiry,
		XRateLimitReset: 30,
		XRateRemaining:  10,
		URL:             body.URL,
		CustomShort:     "",
	}

	database.Client.Decr(database.Ctx, c.ClientIP())

	val, _ = database.Client.Get(database.Ctx, c.ClientIP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := database.Client.TTL(database.Ctx, c.ClientIP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	c.JSON(http.StatusOK, resp)
}
