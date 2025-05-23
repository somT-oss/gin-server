package main

import (
	// "context"
	// "encoding/json"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"github.com/joho/godotenv"

	// "strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
)

type Message struct {
	SourceCode   string `json:"source_code" binding:"required"`
	LanguageType string `json:"language_type" binding:"required"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}


func receiveMessage(c *gin.Context) {
	// Endpoint that receives message from SNS topic
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "The an error occurred"})
		return
	}
	fmt.Println(string(body))
	c.JSON(http.StatusOK, gin.H{"message": string(body)})

	/*
		Functionality for receiving the message from SNS goes here.
		Functionality fo spining up respective docker containers for test goes here as well
	*/
}

func pushMessage(c *gin.Context) {
	// Endpoint that sends message from SQS queue
	if c.ContentType() != "application/json" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"error": "Unsupported media type",
		})
		return
	}
	var message Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON or missing fields",
		})
		return
	}

	// Get AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Unable to load the SDK config, %v", err)
	}

	// Initialize new sqs client
	client := sqs.NewFromConfig(cfg)
	err = godotenv.Load(".env")
	check(err)

	sqsUrl := fmt.Sprintf("%s/%s/%s/", os.Getenv("SQS_URL"), os.Getenv("ACCOUNT_ID"), os.Getenv("QUEUE_NAME"))
	body := fmt.Sprintf("%+v", message)
	resp, err := client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl: &sqsUrl,
		MessageBody: &body,
		DelaySeconds: 0,
	})
	check(err)
	c.JSON(http.StatusOK, gin.H{
		"sqs_message_id": resp.MessageId,
	})
}

/*	- send code through json
	- write code the language specific file
	- create a container to run the code
	- time the container to handle TLE
	- set memory limitations to handle MLE
	- return the response from the code
	- return back to the user
*/


func runContainer() (string, error) {
	cmd := exec.Command(
			"docker",
			"run",
			"-v",
			"/home/somto/Dev/aws-go:/app",
			"-w",
			"/app",
			"python:3",
			"python",
			"testfile.py",
		)
		output, err := cmd.Output()
			if err != nil {
			return "", err
		}
	fmt.Println(string(output))
	return string(output), nil
}

func getCurrentDir() string {
    dir, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }
    return dir
}


func setupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func main() {
	router := setupRouter() 
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	router.POST("/receive", receiveMessage)
	router.POST("/send", pushMessage)

	router.Run("0.0.0.0:8080")
}
