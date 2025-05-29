package main

import (
	"context"
	"flag"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/wh1plash/API/client"
	"github.com/wh1plash/API/service"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log zerolog.Logger

func init() {
	InitLogger("logs/app.log")
	mustLoadEnvVariables()
}

func main() {
	var (
		apiA = flag.String("ApiA", GetVal("GET_URL"), "API A")
		apiB = flag.String("ApiB", GetVal("POST_URL"), "API B")
	)
	flag.Parse()

	//zerolog.SetGlobalLevel(zerolog.InfoLevel)

	Log.Info().Msg("Started")

	retryDelayStr := GetVal("RETRY_DELAY")
	retryDelay, err := time.ParseDuration(retryDelayStr)
	if err != nil {
		Log.Error().Err(err).Str("error to get value for variable", "RETRY_DELAY")
	}

	retryCountStr := GetVal("RETRY_CNT")
	retryCount, err := strconv.Atoi(retryCountStr)
	if err != nil {
		Log.Error().Err(err).Str("error to get value for variable", "RETRY_CNT")
	}

	client := client.NewAPIClient(*apiA, *apiB)

	ctx := context.Background()
	service.DispatchUsers(ctx, client, retryDelay, retryCount)
}

func mustLoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		Log.Error().Err(err).Msg("Error loading .env file.")
	}
}

func InitLogger(logFile string) {
	logWriter := &lumberjack.Logger{
		Filename: logFile,
		MaxSize:  5,
		Compress: true,
	}

	zerolog.TimeFieldFormat = time.RFC3339
	multiWriter := io.MultiWriter(os.Stdout, logWriter)
	Log = zerolog.New(multiWriter).With().Timestamp().Logger()
}

func GetVal(s string) string {
	if os.Getenv(s) == "" {
		Log.Debug().Msg("Setting up defaul value for variable")
		switch s {
		case "RETRY_DELAY":
			return "0.5"
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
