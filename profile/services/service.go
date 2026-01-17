package services

import "date-bot-go/profile"

type ProfileService struct {
	r *profile.Repository
}

func NewProfileService(r *profile.Repository) *ProfileService {
	return &ProfileService{r: r}
}
