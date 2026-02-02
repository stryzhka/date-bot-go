package mock

import (
	"context"
	"date-bot-go/matching/models"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (r *MockRepository) AddLike(ctx context.Context, like *models.Like) error {
	args := r.Called(like)
	return args.Error(0)
}

func (r *MockRepository) DeleteLike(ctx context.Context, like *models.Like) error {
	args := r.Called(like)
	return args.Error(0)
}

func (r *MockRepository) GetUserLikes(ctx context.Context, userId string) []string {
	args := r.Called(userId)
	return args.Get(0).([]string)
}

func (r *MockRepository) IsMutual(ctx context.Context, like *models.Like) (bool, error) {
	args := r.Called(like)
	return args.Bool(0), args.Error(1)
}
