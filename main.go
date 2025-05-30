package main

import (
	"context"
	"flag"
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
	var (
		apiA = flag.String("ApiA", GetVal("GET_URL"), "API A")
		apiB = flag.String("ApiB", GetVal("POST_URL"), "API B")
	)
	flag.Parse()

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

	client := client.NewAPIClient(*apiA, *apiB)

	ctx := context.Background()
	service.DispatchUsers(ctx, client, retryDelay, retryCount)
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
		}
	}
	return os.Getenv(s)
}
