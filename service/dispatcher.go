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

func DispatchUsers(ctx context.Context, client client.UserClient, delay time.Duration, cnt int) {
	users, err := client.GetUsers(ctx)
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}

	for _, user := range users {
		if strings.HasSuffix(user.Email, ".biz") {
			if err := postWithRetry(ctx, client, user, delay, cnt); err != nil {
				log.Printf("Failed to send user %s: %v", user.Email, err)
			} else {
				log.Printf("Sent user %s", user.Email)
			}
		} else {
			log.Printf("Skipping user %s", user.Email)
		}
	}
}

func postWithRetry(ctx context.Context, client client.UserClient, user model.User, delay time.Duration, cnt int) error {
	for i := range cnt {
		err := client.PostUser(ctx, user)
		if err == nil {
			return nil
		}
		log.Printf("Retry %d for user %s", i+1, user.Email)
		time.Sleep(delay)
	}
	return fmt.Errorf("max retries exceeded for user %s", user.Email)
}
