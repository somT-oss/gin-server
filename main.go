package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/hjson/hjson-go"
	"github.com/joho/godotenv"
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

func writeContentToFile(fileContent string, fileType string) {
	if fileType != "python" {
		fmt.Errorf("the file extension must be: 'python'")
		return
	}
	file, err := os.Create("run.py")
	check(err)

	defer file.Close()

	file.Write([]byte(fileContent))
}


func retreiveContainerExitStatus() (string, error) {
	cmd := exec.Command("docker", "wait", "dummy")
	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Convert []byte to string and trim whitespace like \n
	str := strings.TrimSpace(string(res))

	// Parse integer from cleaned string
	exitCode, err := strconv.Atoi(str)
	if err != nil {
		return "", err
	}

	// Return the string version of the integer
	return strconv.Itoa(exitCode), nil
}


func getIntegerContainerExitStatusCode(exitStatus string) (int, error) {
	fmt.Println(exitStatus)
	fmt.Println("Trimmed exit status %s", exitStatus)
	code, err := strconv.Atoi(exitStatus)
	return code, err
}


func runForDuration(ctx context.Context, duration time.Duration, function func()) {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	
	done := make(chan struct{})

	go func() {
		function()
		close(done)
	}()

	select {
		case <- ctx.Done():
			fmt.Println("Time is up")
		case <- done:
			fmt.Println("Function is done")
	}
	
}

func receiveMessage(c *gin.Context) {
	message_id := c.Param("id")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Unable to load the SDK config, %v", err)
	}

	client := sqs.NewFromConfig(cfg)
	sqsUrl := fmt.Sprintf("%s/%s/%s/", os.Getenv("SQS_URL"), os.Getenv("ACCOUNT_ID"), os.Getenv("QUEUE_NAME"))
	resp, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            &sqsUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     10,
		VisibilityTimeout:   5,
	})
	if err != nil {
		log.Fatalf("Failed to receive messages, %v", err)
	}

	if len(resp.Messages) == 0 {
		fmt.Println("No messages received")
	}

	var myMap map[string]interface{}

	for _, message := range resp.Messages {
		if *message.MessageId == message_id {
			err := hjson.Unmarshal([]byte(*message.Body), &myMap)
			check(err)
		}
	}

	if myMap == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not find message in queue",
		})
		return
	}

	mapSourceCode := myMap["source_code"]
	mapLanguageType := myMap["language_type"]
	sourceCode := fmt.Sprintf("%v", mapSourceCode)
	languageType := fmt.Sprintf("%v", mapLanguageType)

	fmt.Println(sourceCode)
	fmt.Println(languageType)

	writeContentToFile(sourceCode, languageType)
	
	context := context.Background()
	duration := 5 * time.Second
	runForDuration(context, duration, runContainer)
	
	containerStatusCode, err := retreiveContainerExitStatus()
	check(err)
	
	actualStatusCode, err := getIntegerContainerExitStatusCode(containerStatusCode)
	check(err)

	if actualStatusCode == 1 {
		log.Printf("The container did not exit gracefully")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The container did not exit gracefully",
		})
		return
	}	else if actualStatusCode == 137 {
		log.Printf("Memory Limit Exceeded")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Memory Limit Exceeded",
		})
		return
	}

	if actualStatusCode != 0 {
		log.Printf("An unrecognized error occurred %v", actualStatusCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "An unrecognized error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your code passed all the test cases.",
	})
}

func pushMessage(c *gin.Context) {
	// Endpoint that sends message from SQS queue
	if c.ContentType() != "application/json" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"error": "Unsupported media type",
		})
		return
	}

	requestBodyBytes, _ := io.ReadAll(c.Request.Body)
	requestBody := string(requestBodyBytes)

	// TODO: Check how to confirm request body format and still get the request body after processing.

	// var message Message
	// if err := c.ShouldBindJSON(&message); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Invalid JSON or missing fields",
	// 	})
	// 	return
	// }

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
	resp, err := client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:     &sqsUrl,
		MessageBody:  &requestBody,
		DelaySeconds: 0,
	})
	check(err)
	c.JSON(http.StatusOK, gin.H{
		"request_body":   requestBody,
		"sqs_message_id": resp.MessageId,
	})
}

/*	- send code through json -- done
	- write code the language specific file
	- create a container to run the code
	- time the container to handle TLE
	- set memory limitations to handle MLE
	- return the response from the code
	- return back to the user
*/

func runContainer() {
	cmd := exec.Command(
		"docker",
		"run",
		"--name",
		"dummy",
		"-d",
		"-v",
		fmt.Sprintf("%v:/app", getCurrentDir()),
		"-w",
		"/app",
		"python:3",
		"python",
		"test.py",
	)
	container_id, err := cmd.Output()
	check(err)
	fmt.Println(string(container_id))
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
	router.GET("/receive/:id", receiveMessage)
	router.POST("/send", pushMessage)

	router.Run("0.0.0.0:8080")
}
