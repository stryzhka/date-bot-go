package mock

import (
	"context"
	"date-bot-go/profile/models"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (s *MockService) Create(
	ctx context.Context, userId, name, gender, description string,
) error {
	args := s.Called(ctx, userId, name, gender, description)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (s *MockService) GetById(ctx context.Context, id string) *models.Profile {
	args := s.Called(id)
	if args.Get(0) == nil {
		return nil
	}
	//TODO: видимо это и в других тестах
	return args.Get(0).(*models.Profile)
}

func (s *MockService) GetAll(ctx context.Context) []models.Profile {
	args := s.Called()
	return args.Get(0).([]models.Profile)
}

func (s *MockService) UpdateById(ctx context.Context, id string, newProfile *models.Profile) error {
	args := s.Called(ctx, id, newProfile)
	return args.Error(0)
}

func (s *MockService) DeleteById(ctx context.Context, id string) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}
