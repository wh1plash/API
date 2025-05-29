package main

import (
	"context"
	"flag"
	"log"
	"os"
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
		cnt  int
		apiA = flag.String("ApiA", os.Getenv("GET_URL"), "API A")
		apiB = flag.String("ApiB", os.Getenv("POST_URL"), "API B")
	)
	flag.Parse()

	retryDelayStr := os.Getenv("RETRY_DELAY")
	if retryDelayStr == "" {
		retryDelayStr = "2s"
	}
	retryDelay, err := time.ParseDuration(retryDelayStr)
	if err != nil {
		log.Fatalf("Invalid delay value: %v", err)
	}

	retryCnt := os.Getenv("RETRY_CNT")
	if retryCnt == "" {
		cnt = 3
	}

	client := client.NewAPIClient(*apiA, *apiB)

	ctx := context.Background()
	service.DispatchUsers(ctx, client, retryDelay, cnt)
}

func mustLoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}
}
