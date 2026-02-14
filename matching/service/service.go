package service

import (
	"context"
	"date-bot-go/matching"
	"date-bot-go/matching/client"
	"date-bot-go/matching/models"
	"log"
	"math/rand/v2"
	"slices"
)

type MatchingService struct {
	r               matching.Repository
	profileProvider client.ProfileProvider
}

func NewMatchingService(r matching.Repository, profileProvider client.ProfileProvider) *MatchingService {
	return &MatchingService{r: r, profileProvider: profileProvider}
}

func (s *MatchingService) Like(ctx context.Context, userId, likedId string) error {
	if userId == likedId {
		return matching.ErrAutoLike
	}
	like := &models.Like{
		UserId:  userId,
		LikedId: likedId,
	}
	like1 := &models.Like{
		UserId:  likedId,
		LikedId: userId,
	}

	//TODO: колхоз!!!
	err := s.r.AddLike(ctx, like)
	mutual, err := s.r.IsMutual(ctx, like)
	if err != nil {
		return err
	}
	if mutual {
		//--send link to likedId
		//--send link to userId
		log.Println("stab: mutual match: ", userId, ", ", likedId)
		err = s.r.DeleteLike(ctx, like)
		err = s.r.DeleteLike(ctx, like1)
		return err
	}
	log.Println("stab: ", likedId, " got like from ", userId)
	//--notify likedId!
	//--next profile
	return err
}

// TODO: другой алгоримт
func (s *MatchingService) NextProfile(ctx context.Context, userId string) (*models.Profile, error) {
	userLikes := s.r.GetUserLikes(ctx, userId)
	users, err := s.profileProvider.GetCandidates(ctx)
	var filtered []models.Profile
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if !slices.Contains(userLikes, user.UserId) && user.UserId != userId {
			filtered = append(filtered, user)
		}
	}
	l := len(filtered)
	index := l
	if l > 0 {
		index = rand.IntN(l)
		return &filtered[index], nil
	}
	//но так не должно быть
	return nil, nil
}
