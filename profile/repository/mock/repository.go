package mock

import (
	"context"
	"date-bot-go/profile/models"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (r *MockRepository) Create(ctx context.Context, profile *models.Profile) error {
	args := r.Called(profile)
	return args.Error(0)
}

func (r *MockRepository) Get(ctx context.Context, id string) *models.Profile {
	args := r.Called(id)
	return args.Get(0).(*models.Profile)
}

func (r *MockRepository) GetAll(ctx context.Context) []models.Profile {
	args := r.Called()
	return args.Get(0).([]models.Profile)
}

func (r *MockRepository) UpdateById(ctx context.Context, id string, newProfile *models.Profile) error {
	args := r.Called(id, profile)
	return args.Error(0)
}

func (r *MockRepository) DeleteById(ctx context.Context, id string) error {
	args := r.Called(id)
	return args.Error(0)
}
