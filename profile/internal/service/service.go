package service

import (
	"context"
	"github.com/google/uuid"
	profile "profile/internal"
	"profile/internal/models"
	"strings"
	"time"
)

type ProfileService struct {
	r profile.Repository
}

func NewProfileService(r profile.Repository) *ProfileService {
	return &ProfileService{r: r}
}

func (p *ProfileService) Create(
	ctx context.Context, userId, name, gender, description string,
) error {
	//if strings.TrimSpace(userId) == "" {
	//	return profile.ErrValidationUserId
	//}
	if strings.TrimSpace(name) == "" || len(name) >= 15 {
		return profile.ErrValidationName
	}
	profile := &models.Profile{
		Id:          uuid.New(),
		UserId:      userId,
		Name:        name,
		Gender:      gender,
		Description: description,
		Topics:      nil,
		DateCreated: time.Now(),
		PhotoPath:   "",
	}
	return p.r.Create(ctx, profile)
}

func (p *ProfileService) GetById(ctx context.Context, id string) *models.Profile {
	return p.r.Get(ctx, id)
}

func (p *ProfileService) GetAll(ctx context.Context) []models.Profile {
	return p.r.GetAll(ctx)
}

func (p *ProfileService) UpdateById(ctx context.Context, id string, newProfile *models.Profile) error {
	if strings.TrimSpace(newProfile.Name) == "" || len(newProfile.Name) >= 15 {
		return profile.ErrValidationName
	}
	return p.r.UpdateById(ctx, id, newProfile)
}

func (p *ProfileService) DeleteById(ctx context.Context, id string) error {
	return p.r.DeleteById(ctx, id)
}
