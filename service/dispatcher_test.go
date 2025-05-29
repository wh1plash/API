package service

import (
	"context"
	"testing"
	"time"

	"github.com/wh1plash/API/client"
	"github.com/wh1plash/API/model"
)

func TestEmail(t *testing.T) {
	mock := &client.MockClient{
		UsersToReturn: []model.User{
			{Name: "Alice", Email: "alice@corp.biz"},
			{Name: "Bob", Email: "bob@example.com"},
			{Name: "Eve", Email: "eve@evil.biz"},
		},
	}

	ctx := context.Background()
	DispatchUsers(ctx, mock, 0, 3)

	expectedEmails := []string{"alice@corp.biz", "eve@evil.biz"}

	if len(expectedEmails) != len(mock.PostedUsers) {
		t.Fatalf("expected %d posted users, but got %d", len(expectedEmails), len(mock.PostedUsers))
	}

	for i, user := range mock.PostedUsers {
		if user.Email != expectedEmails[i] {
			t.Errorf("unexpected user posted at index %d: got %s, want %s", i, user.Email, expectedEmails[i])
		}
	}
}

func TestRetriesWithSuccess(t *testing.T) {
	mock := &client.MockClient{
		UsersToReturn: []model.User{
			{Name: "RetryUser", Email: "fail@retry.biz"},
		},
		PostFailCount: 3,
	}

	ctx := context.Background()
	DispatchUsers(ctx, mock, 1*time.Second, 4)

	if len(mock.PostedUsers) != 1 {
		t.Fatalf("expected 1 posted user, got %d", len(mock.PostedUsers))
	}
	if mock.PostedUsers[0].Email != "fail@retry.biz" {
		t.Errorf("expected user email fail@retry.biz, got %s", mock.PostedUsers[0].Email)
	}

	if mock.PostFailCount >= mock.PostCalls {
		t.Errorf("expected retry to eventually succeed, got %d calls", mock.PostCalls)
	}
}
