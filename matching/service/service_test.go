package service

import (
	"context"
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
