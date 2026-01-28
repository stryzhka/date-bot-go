package profile

import (
	"context"
	"date-bot-go/profile/models"
)

type Repository interface {
	Create(ctx context.Context, profile *models.Profile) error
	Get(ctx context.Context, id string) *models.Profile
	GetAll(ctx context.Context) []models.Profile
	UpdateById(ctx context.Context, id string, newProfile *models.Profile) error
	DeleteById(ctx context.Context, id string) error
	// TODO AddTopicById
	// TODO CleanTopicsById
}
