package service

import (
	"context"
	"date-bot-go/matching"
	"date-bot-go/matching/client"
	"date-bot-go/matching/models"
)

type MatchingService struct {
	r               matching.Repository
	profileProvider client.ProfileProvider
}

func NewMatchingService(r matching.Repository, profileProvider client.ProfileProvider) *MatchingService {
	return &MatchingService{r: r, profileProvider: profileProvider}
}

func (s *MatchingService) Like(ctx context.Context, userId, likedId string) error {
	like := &models.Like{
		UserId:  userId,
		LikedId: likedId,
	}
	like1 := &models.Like{
		UserId:  likedId,
		LikedId: userId,
	}
	mutual, err := s.r.IsMutual(ctx, like)
	if err != nil {
		return err
	}
	if mutual {
		//--send link to likedId
		//--send link to userId
		err = s.r.DeleteLike(ctx, like)
		err = s.r.DeleteLike(ctx, like1)
		return err
	}
	err = s.r.AddLike(ctx, like)
	//--notify likedId!
	//--next profile
	return err
}

func (s *MatchingService) NextProfile(ctx context.Context, userId string) (*models.Profile, error) {
	//???
	return s.profileProvider.GetCandidate(ctx, userId)
}
