package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wh1plash/API/model"
)

func DispatchUsers(ctx context.Context, users []model.User, postURL string) {
	for _, user := range users {
		if strings.HasSuffix(user.Email, ".biz") {
			if err := postWithRetry(ctx, postURL, user); err != nil {
				log.Printf("Failed to send user %s after retries: %v\n", user.Email, err)
			} else {
				log.Printf("Sending user %s to API B", user.Name)
			}
			log.Printf("Skipping user %s (not .biz)", user.Email)
		}

	}
}

func postWithRetry(ctx context.Context, url string, user model.User) error {
	body := map[string]string{
		"name":  user.Name,
		"email": user.Email,
	}
	data, _ := json.Marshal(body)

	for i := 1; i <= 3; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
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
		log.Printf("Retry %d for user %s\n", i, user.Email)
		delayStr := os.Getenv("RETRY_DELAY")
		if delayStr == "" {
			delayStr = "1s" // default
		}
		delay, _ := time.ParseDuration(os.Getenv("DELAY"))

		time.Sleep(time.Second * delay)
	}
	return fmt.Errorf("max retries exceeded for user %s", user.Email)

}
