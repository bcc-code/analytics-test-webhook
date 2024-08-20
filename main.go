package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type DataStore struct {
	sync.Mutex
	data          map[string][]json.RawMessage
	lastUpdatedAt map[string]time.Time
}

var (
	store = DataStore{
		data:          make(map[string][]json.RawMessage),
		lastUpdatedAt: make(map[string]time.Time),
	}
	apiKey = os.Getenv("API_KEY")
)

func webhookHandler(c *gin.Context) {
	requestApiKey := c.Query("api_key")
	if requestApiKey != apiKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID"})
		return
	}

	store.Lock()
	defer store.Unlock()

	if lastUpdate, exists := store.lastUpdatedAt[id]; exists && time.Since(lastUpdate) > 10*time.Minute {
		store.data[id] = nil
	}

	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	store.data[id] = append(store.data[id], payload)
	store.lastUpdatedAt[id] = time.Now()

	c.Status(http.StatusOK)
}

func getDataHandler(c *gin.Context) {
	requestApiKey := c.Query("api_key")
	if requestApiKey != apiKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID"})
		return
	}

	store.Lock()
	defer store.Unlock()

	if lastUpdate, exists := store.lastUpdatedAt[id]; exists && time.Since(lastUpdate) > 10*time.Minute {
		store.data[id] = nil
	}

	data, exists := store.data[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func main() {
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is not set")
	}

	router := gin.Default()

	router.POST("/webhook/:id", webhookHandler)
	router.GET("/get_data/:id", getDataHandler)

	log.Println("Server is running on port 8080")
	log.Fatal(router.Run(":8080"))
}
