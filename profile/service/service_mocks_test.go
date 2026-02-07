package service

import (
	"context"
	"date-bot-go/profile"
	"date-bot-go/profile/models"
	"date-bot-go/profile/repository/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestSuccessCreate(t *testing.T) {
	profileRepository := new(mock.MockRepository)
	profileService := NewProfileService(profileRepository)
	var (
		ctx         = context.Background()
		userId      = "123"
		name        = "test"
		gender      = "f"
		description = "test test test"
	)
	profileRepository.On("Create", mock2.Anything).Return(nil)
	err := profileService.Create(ctx, userId, name, gender, description)
	assert.NoError(t, err)
}

func TestUserAlreadyExistsCreate(t *testing.T) {
	profileRepository := new(mock.MockRepository)
	profileService := NewProfileService(profileRepository)
	var (
		ctx         = context.Background()
		userId      = "123"
		name        = "test"
		gender      = "f"
		description = "test test test"
	)
	profileRepository.On("Create", mock2.Anything).Return(profile.ErrUserAlreadyExists)
	err := profileService.Create(ctx, userId, name, gender, description)
	assert.Error(t, err)
}

func TestFailCreate(t *testing.T) {
	profileRepository := new(mock.MockRepository)
	profileService := NewProfileService(profileRepository)
	var (
		ctx         = context.Background()
		userId      = "     "
		name        = "sfdanlnjladsjkblgdkjblgdakjjda"
		gender      = "f"
		description = "test test test"
	)
	profileRepository.On("Create", mock2.Anything).Return(nil)
	err := profileService.Create(ctx, userId, name, gender, description)
	assert.Error(t, err)
}

func TestSuccessUpdateById(t *testing.T) {
	profileRepository := new(mock.MockRepository)
	profileService := NewProfileService(profileRepository)
	var (
		ctx        = context.Background()
		id         = "123"
		newProfile = models.Profile{
			Id:          uuid.New(),
			UserId:      "123",
			Name:        "test",
			Gender:      "f",
			Description: "test test test",
			Topics:      nil,
			DateCreated: time.Now(),
			PhotoPath:   "",
		}
	)
	profileRepository.On("UpdateById", id, mock2.Anything).Return(nil)
	err := profileService.UpdateById(ctx, id, &newProfile)
	assert.NoError(t, err)
}

func TestFailUpdateById(t *testing.T) {
	profileRepository := new(mock.MockRepository)
	profileService := NewProfileService(profileRepository)
	var (
		ctx        = context.Background()
		id         = "123"
		newProfile = models.Profile{
			Id:          uuid.New(),
			UserId:      "123",
			Name:        "tttttttttttttttttttttttttttttttttttttttt",
			Gender:      "f",
			Description: "test test test",
			Topics:      nil,
			DateCreated: time.Now(),
			PhotoPath:   "",
		}
	)
	profileRepository.On("UpdateById", id, mock2.Anything).Return(nil)
	err := profileService.UpdateById(ctx, id, &newProfile)
	assert.Error(t, err)
}

func TestUserNotFoundUpdateById(t *testing.T) {
	profileRepository := new(mock.MockRepository)
	profileService := NewProfileService(profileRepository)
	var (
		ctx        = context.Background()
		id         = "123"
		newProfile = models.Profile{
			Id:          uuid.New(),
			UserId:      "123",
			Name:        "tttttttttttttttttttttttttttttttttttttttt",
			Gender:      "f",
			Description: "test test test",
			Topics:      nil,
			DateCreated: time.Now(),
			PhotoPath:   "",
		}
	)
	profileRepository.On("UpdateById", id, mock2.Anything).Return(profile.ErrUserNotFound)
	err := profileService.UpdateById(ctx, id, &newProfile)
	assert.Error(t, err)
}

func TestUserNotFoundDeleteById(t *testing.T) {
	profileRepository := new(mock.MockRepository)
	profileService := NewProfileService(profileRepository)
	var (
		ctx = context.Background()
		id  = "123"
	)
	profileRepository.On("DeleteById", id, mock2.Anything).Return(profile.ErrUserNotFound)
	err := profileService.DeleteById(ctx, id)
	assert.Error(t, err)
}
