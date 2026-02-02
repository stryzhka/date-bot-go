package service

import (
	"context"
	"date-bot-go/matching"
	mockProvider "date-bot-go/matching/client/mock"
	"date-bot-go/matching/models"
	"date-bot-go/matching/repository/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSuccessNoMutualLike(t *testing.T) {
	r := new(mock.MockRepository)
	profileProvider := new(mockProvider.MockProfileProvider)
	s := NewMatchingService(r, profileProvider)
	expectedLike := &models.Like{
		UserId:  "123",
		LikedId: "456",
	}
	r.On("AddLike", expectedLike).Return(nil)
	r.On("IsMutual", expectedLike).Return(false, nil)
	err := s.Like(context.Background(), "123", "456")
	assert.NoError(t, err)
}

// for case when 456 liked 123 before
func TestSuccessMutualLike(t *testing.T) {
	r := new(mock.MockRepository)
	profileProvider := new(mockProvider.MockProfileProvider)
	s := NewMatchingService(r, profileProvider)
	expectedLike := &models.Like{
		UserId:  "123",
		LikedId: "456",
	}
	r.On("AddLike", expectedLike).Return(nil)
	r.On("DeleteLike", expectedLike.UserId).Return(nil)
	r.On("DeleteLike", expectedLike.LikedId).Return(nil)
	r.On("IsMutual", expectedLike).Return(true, nil)
	err := s.Like(context.Background(), "123", "456")
	assert.NoError(t, err)
}

// calling next profile from separate handler
func TestNextProfile(t *testing.T) {
	r := new(mock.MockRepository)
	profileProvider := new(mockProvider.MockProfileProvider)
	userId := "123"
	profileProvider.On("GetCandidate", userId).Return(&models.Profile{}, nil)
	s := NewMatchingService(r, profileProvider)
	expectedProfile, err := s.NextProfile(context.Background(), userId)
	assert.NoError(t, err)
	assert.IsType(t, expectedProfile, &models.Profile{})
}

func TestErrAutoLikeAddLike(t *testing.T) {
	r := new(mock.MockRepository)
	profileProvider := new(mockProvider.MockProfileProvider)
	s := NewMatchingService(r, profileProvider)
	expectedLike := &models.Like{
		UserId:  "123",
		LikedId: "123",
	}
	r.On("AddLike", expectedLike).Return(nil)
	err := s.Like(context.Background(), "123", "123")
	assert.Error(t, err)
	assert.IsType(t, matching.ErrAutoLike, err)
}

func TestSuccessNoLikesNextProfile(t *testing.T) {
	//исключить только вызывающий id
	r := new(mock.MockRepository)
	profileProvider := new(mockProvider.MockProfileProvider)
	s := NewMatchingService(r, profileProvider)
	r.On("GetUserLikes", "123").Return([]string{}, nil)
	var allProfiles = []models.Profile{
		{
			UserId:      "123",
			Username:    "test",
			Name:        "test",
			Gender:      "f",
			Description: "test test test",
			PhotoPath:   "test",
		},
		{
			UserId:      "456",
			Username:    "test1",
			Name:        "test1",
			Gender:      "f",
			Description: "test test test",
			PhotoPath:   "test",
		},
		{
			UserId:      "789",
			Username:    "test2",
			Name:        "test2",
			Gender:      "f",
			Description: "test test test",
			PhotoPath:   "test",
		},
	}
	profileProvider.On("GetCandidates").Return(allProfiles, nil)
	expProfile, err := s.NextProfile(context.Background(), "123")
	assert.NoError(t, err)
	assert.IsType(t, expProfile, &models.Profile{})
	t.Log("Next profile: ", expProfile)
}

func TestSuccessHasLikesNextProfile(t *testing.T) {
	r := new(mock.MockRepository)
	profileProvider := new(mockProvider.MockProfileProvider)
	s := NewMatchingService(r, profileProvider)
	r.On("GetUserLikes", "123").Return([]string{"456", "789"}, nil)
	var allProfiles = []models.Profile{
		{
			UserId:      "123",
			Username:    "test",
			Name:        "test",
			Gender:      "f",
			Description: "test test test",
			PhotoPath:   "test",
		},
		{
			UserId:      "456",
			Username:    "test1",
			Name:        "test1",
			Gender:      "f",
			Description: "test test test",
			PhotoPath:   "test",
		},
		{
			UserId:      "789",
			Username:    "test2",
			Name:        "test2",
			Gender:      "f",
			Description: "test test test",
			PhotoPath:   "test",
		},
	}
	profileProvider.On("GetCandidates").Return(allProfiles, nil)
	expProfile, err := s.NextProfile(context.Background(), "123")
	assert.NoError(t, err)
	assert.IsType(t, expProfile, &models.Profile{})
	//но так не должно быть
	t.Log("Next profile: ", expProfile)
}
