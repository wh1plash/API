package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wh1plash/API/model"
)

func FetchUser(ctx context.Context, url string) ([]model.User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
