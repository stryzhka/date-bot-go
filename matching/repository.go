package matching

import (
	"context"
	"date-bot-go/matching/models"
)

type Repository interface {
	AddLike(ctx context.Context, like *models.Like) error
	DeleteLike(ctx context.Context, userId string) error
	IsMutual(ctx context.Context, like *models.Like) (bool, error)
}
