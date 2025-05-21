package main

import (
	// "context"
	// "encoding/json"
	"fmt"
	"net/http"

	// "strconv"

	"github.com/gin-gonic/gin"
	// "fmt"
	// "log"
	// "os"
	// "github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go-v2/service/sqs"
)

type MyPayload struct {
	Name string `json:"name" binding:"required"`
}

type Album struct {
	ID       int      `json:"id" binding:"required"`
	Title    string   `json:"title" binding:"required"`
	Artistes []string `json:"artistes" binding:"required"`
}

var albums []Album

func addAlbum(c *gin.Context) {
	if c.ContentType() != "application/json" {
		c.JSON(405, "Unsupported content type")
		return
	}
	var album Album
	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
	}
	albums = append(albums, album)
	c.JSON(http.StatusOK, albums)
}

func printPostRequest(c *gin.Context) {
	if c.ContentType() != "application/json" {
		c.JSON(405, "Unsupported content type")
		return
	}

	var payload MyPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON or missing fields"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successful",
		"name":    payload.Name,
	})
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/response", printPostRequest)
	router.POST("/add-album", addAlbum)

	fmt.Println("running server on localhost:8000")
	router.Run("0.0.0.0:8080")
}
