package mock

import (
	"context"
	"date-bot-go/matching/models"
	"github.com/stretchr/testify/mock"
)

type MockProfileProvider struct {
	mock.Mock
}

func (p *MockProfileProvider) GetCandidate(ctx context.Context, excludeId string) *models.Profile {
	args := p.Called(excludeId)
	return args.Get(0).(*models.Profile)
}
