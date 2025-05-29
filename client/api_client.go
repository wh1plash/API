package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wh1plash/API/model"
)

type UserClient interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	PostUser(ctx context.Context, user model.User) error
}

type APIClient struct {
	apiA string
	apiB string
}

func NewAPIClient(apiA, apiB string) *APIClient {
	return &APIClient{
		apiA: apiA,
		apiB: apiB,
	}
}

func (c *APIClient) GetUsers(ctx context.Context) ([]model.User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiA, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	var users []model.User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *APIClient) PostUser(ctx context.Context, user model.User) error {
	body := map[string]string{
		"name":  user.Name,
		"email": user.Email,
	}
	data, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiB, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		resp.Body.Close()
		return nil
	}
	if resp != nil {
		resp.Body.Close()
	}
	return fmt.Errorf("failed to POST user: %s", user.Email)
}
