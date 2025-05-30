package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/wh1plash/API/client"
	"github.com/wh1plash/API/service"
)

func init() {
	mustLoadEnvVariables()
}

func main() {
	apiA := GetVal("GET_URL")
	apiB := GetVal("POST_URL")

	retryDelayStr := GetVal("RETRY_DELAY")
	retryDelay, err := time.ParseDuration(retryDelayStr)
	if err != nil {
		log.Fatal("error to get value for variable", err)
	}

	retryCountStr := GetVal("RETRY_CNT")
	retryCount, err := strconv.Atoi(retryCountStr)
	if err != nil {
		log.Fatal("error to get value for variable", err)
	}

	timeoutStr := GetVal("POST_TIMEOUT")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatal("error to get value for variable", err)
	}

	client := client.NewAPIClient(apiA, apiB)

	service.DispatchUsers(client, retryDelay, retryCount, timeout)
}

func mustLoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func GetVal(s string) string {
	if os.Getenv(s) == "" {
		switch s {
		case "RETRY_DELAY":
			return "0.5s"
		case "GET_URL":
			return "https://jsonplaceholder.typicode.com/users"
		case "POST_URL":
			return "https://webhook.site"
		case "RETRY_CNT":
			return "3"
		case "POST_TIMEOUT":
			return "100ms"
		}
	}
	return os.Getenv(s)
}
