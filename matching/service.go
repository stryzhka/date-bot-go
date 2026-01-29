package matching

import (
	"context"
	"date-bot-go/matching/models"
)

type Service interface {
	Like(ctx context.Context, userId, likedId string) error
	NextProfile(ctx context.Context) *models.Profile
}
