package grpc

import (
	"context"
	"date-bot-go/matching/models"
	"date-bot-go/pkg/profilepb"
	"google.golang.org/grpc"
)

type GrpcProfileProvider struct {
	client profilepb.ProfileServiceClient
	conn   *grpc.ClientConn
}

func NewGrpcProfileProvider(client profilepb.ProfileServiceClient, conn *grpc.ClientConn) *GrpcProfileProvider {
	return &GrpcProfileProvider{
		client: client,
		conn:   conn,
	}
}

func (p *GrpcProfileProvider) GetCandidate(ctx context.Context, excludeId string) (*models.Profile, error) {
	req := &profilepb.GetAllRequest{}
	resp, err := p.client.GetAll(ctx, req)
	profile := &models.Profile{
		UserId:      resp.Profiles[0].UserId,
		Username:    "returned",
		Name:        resp.Profiles[0].Name,
		Gender:      resp.Profiles[0].Gender,
		Description: resp.Profiles[0].Description,
		PhotoPath:   resp.Profiles[0].PhotoPath,
	}
	return profile, err
}
