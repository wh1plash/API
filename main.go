package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wh1plash/API/client"
)

func init() {
	mustLoadEnvVariables()
}

func main() {
	var (
		src  = flag.String("A", GetUrl(), "API A")
		dest = flag.String("B", PostUrl(), "API B")
	)
	flag.Parse()

	ctx := context.Background()
	users, err := client.FetchUser(*src)
	if err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}

	fmt.Println(users)
	_ = ctx
	_ = dest
}

func mustLoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}
}

func GetUrl() string {
	return os.Getenv("GET_URL")
}

func PostUrl() string {
	return os.Getenv("POST_URL")
}
