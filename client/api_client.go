package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/wh1plash/API/model"
	"gopkg.in/natefinch/lumberjack.v2"
)

type UserClient interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	PostUser(ctx context.Context, user model.User) error
}

type APIClient struct {
	apiA   string
	apiB   string
	logger zerolog.Logger
}

func NewAPIClient(apiA, apiB string) *APIClient {
	return &APIClient{
		apiA:   apiA,
		apiB:   apiB,
		logger: InitLogger("logs/app.log"),
	}
}

func (c *APIClient) GetUsers(ctx context.Context) ([]model.User, error) {
	start := time.Now()
	c.logger.Info().Str("ulr", c.apiA).Msg("Fetching users")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiA, nil)
	if err != nil {
		c.logger.Error().Err(err).Str("url", c.apiA).Msg("Failed to create GET request")
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Error().Err(err).Str("url", c.apiA).Dur("duration", time.Since(start)).Msg("GET request failed")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	var users []model.User
	if err := json.Unmarshal(body, &users); err != nil {
		c.logger.Error().Err(err).Msg("Failed to decode user response")
		return nil, err
	}
	c.logger.Info().Dur("duration", time.Since(start)).Msg(fmt.Sprintf("Successfully getting %d users", len(users)))
	return users, nil

}

func (c *APIClient) PostUser(ctx context.Context, user model.User) error {
	start := time.Now()
	body := map[string]string{
		"name":  user.Name,
		"email": user.Email,
	}
	data, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiB, bytes.NewBuffer(data))
	if err != nil {
		c.logger.Error().Err(err).Str("url", c.apiB).Msg("Failed to create POST request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Error().
			Err(err).
			Str("email", user.Email).
			Str("url", c.apiB).
			Dur("duration", time.Since(start)).
			Msg("POST request failed")
		return err
	}
	defer resp.Body.Close()

	c.logger.Info().
		Str("email", user.Email).
		Str("method", http.MethodPost).
		Int("status", resp.StatusCode).
		Dur("duration", time.Since(start)).
		Msg("POST user")

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("POST failed: status %d", resp.StatusCode)
	}
	if resp != nil {
		resp.Body.Close()
	}
	return nil

}

func InitLogger(logFile string) zerolog.Logger {
	logWriter := &lumberjack.Logger{
		Filename: logFile,
		MaxSize:  5,
		Compress: true,
	}

	zerolog.TimeFieldFormat = time.RFC3339
	//multiWriter := io.MultiWriter(os.Stdout, logWriter)
	return zerolog.New(logWriter).With().Timestamp().Logger()
}
