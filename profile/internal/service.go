package profile

import (
	"context"
	"profile/internal/models"
)

type Service interface {
	Create(
		ctx context.Context, userId, name, gender, description string,
	) error
	GetById(ctx context.Context, id string) *models.Profile
	GetAll(ctx context.Context) []models.Profile
	UpdateById(ctx context.Context, id string, newProfile *models.Profile) error
	DeleteById(ctx context.Context, id string) error
}
