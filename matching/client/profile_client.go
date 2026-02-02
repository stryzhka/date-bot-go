package client

import (
	"context"
	"date-bot-go/matching/models"
)

type ProfileProvider interface {
	//TODO all...
	GetCandidates(ctx context.Context) ([]models.Profile, error)
}
