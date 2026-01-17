package profile

import "date-bot-go/profile/models"

type Service interface {
	CreateProfile(profile *models.Profile) error
	GetProfileById(id string) (*models.Profile, string)
	UpdateProfileById(id string, profile *models.Profile) error
	DeleteProfileById(id string) error
}
