package matching

import (
	"context"
	"date-bot-go/matching/models"
)

type Service interface {
	Like(ctx context.Context, userId, likedId string) error
	//GetUserLikes(ctx context.Context, userId string)
	NextProfile(ctx context.Context) (*models.Profile, error)
}
