package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable not set")
	}

	// Create a new Gin router
	r := gin.Default()

	// Define routes
	r.GET("/sprinkleNow/:lat/:lon", func(c *gin.Context) {
		lat, err := strconv.ParseFloat(c.Param("lat"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
			return
		}

		lon, err := strconv.ParseFloat(c.Param("lon"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
			return
		}

		// Get optional soil moisture
		var soilMoisture *float64
		soilMoistureStr := c.Query("soilm")
		if soilMoistureStr != "" {
			sm, err := strconv.ParseFloat(soilMoistureStr, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid soil moisture"})
				return
			}
			soilMoisture = &sm
		}

		// Create sprinkler instance
		sprinkler := NewOpenWeatherSprinkler(apiKey, lat, lon, nil, nil, soilMoisture)

		// Return result
		c.JSON(http.StatusOK, gin.H{"sprinkle": sprinkler.SprinkleNow()})
	})

	r.GET("/sprinkleTime/:lat/:lon/:barrelv/:pumpo", func(c *gin.Context) {
		lat, err := strconv.ParseFloat(c.Param("lat"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
			return
		}

		lon, err := strconv.ParseFloat(c.Param("lon"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
			return
		}

		barrelVolume, err := strconv.ParseFloat(c.Param("barrelv"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid barrel volume"})
			return
		}

		pumpOutput, err := strconv.ParseFloat(c.Param("pumpo"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pump output"})
			return
		}

		// Get optional soil moisture
		var soilMoisture *float64
		soilMoistureStr := c.Query("soilm")
		if soilMoistureStr != "" {
			sm, err := strconv.ParseFloat(soilMoistureStr, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid soil moisture"})
				return
			}
			soilMoisture = &sm
		}

		// Create sprinkler instance
		sprinkler := NewOpenWeatherSprinkler(apiKey, lat, lon, &barrelVolume, &pumpOutput, soilMoisture)

		// Return result
		c.JSON(http.StatusOK, gin.H{"sprinkleTime": sprinkler.GetSprinkleTime()})
	})

	// Run the server
	r.Run(":8080")
}
