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

func (p *GrpcProfileProvider) GetCandidates(ctx context.Context) ([]models.Profile, error) {
	req := &profilepb.GetAllRequest{}
	resp, err := p.client.GetAll(ctx, req)
	var profiles []models.Profile = nil
	for _, val := range resp.Profiles {
		profile := &models.Profile{
			UserId:      val.UserId,
			Username:    val.Username,
			Name:        val.Name,
			Gender:      val.Gender,
			Description: val.Description,
			PhotoPath:   val.PhotoPath,
		}
		profiles = append(profiles, *profile)
	}
	return profiles, err
}
