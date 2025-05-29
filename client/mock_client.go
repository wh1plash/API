package client

import (
	"context"
	"errors"

	"github.com/wh1plash/API/model"
)

type MockClient struct {
	UsersToReturn []model.User
	PostFailCount int
	PostedUsers   []model.User
	PostCalls     int
}

func (m *MockClient) GetUsers(ctx context.Context) ([]model.User, error) {
	return m.UsersToReturn, nil
}

func (m *MockClient) PostUser(ctx context.Context, user model.User) error {
	m.PostCalls++
	if m.PostCalls <= m.PostFailCount {
		return errors.New("simulated failure")
	}
	m.PostedUsers = append(m.PostedUsers, user)
	return nil
}
