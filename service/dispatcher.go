package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/wh1plash/API/client"
	"github.com/wh1plash/API/model"
)

func DispatchUsers(client client.UserClient, delay time.Duration, cnt int, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	users, err := client.GetUsers(ctx)
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
		return
	}
	cancel()

	for _, user := range users {
		if strings.HasSuffix(user.Email, ".biz") {
			if err := postWithRetry(client, user, delay, cnt, timeout); err != nil {
				log.Printf("Failed to send user %s: %v", user.Email, err)
			}
		} else {
			log.Printf("Skipping user %s", user.Email)
		}
	}
}

func postWithRetry(client client.UserClient, user model.User, delay time.Duration, cnt int, timeout time.Duration) error {
	for i := range cnt {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err := client.PostUser(ctx, user)
		if err == nil {
			return nil
		}
		log.Printf("Retry %d for user %s", i+1, user.Email)
		time.Sleep(delay)
	}
	return fmt.Errorf("max retries exceeded for user %s", user.Email)
}
