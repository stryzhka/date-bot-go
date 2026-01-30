package client

import (
	"context"
	"date-bot-go/matching/models"
)

type ProfileProvider interface {
	//TODO all...
	GetCandidate(ctx context.Context, excludeId string) (*models.Profile, error)
}
