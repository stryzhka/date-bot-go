package mock

import (
	"context"
	"date-bot-go/matching/models"
	"github.com/stretchr/testify/mock"
)

type MockProfileProvider struct {
	mock.Mock
}

func (p *MockProfileProvider) GetCandidates(ctx context.Context) ([]models.Profile, error) {
	args := p.Called()
	return args.Get(0).([]models.Profile), args.Error(1)
}
