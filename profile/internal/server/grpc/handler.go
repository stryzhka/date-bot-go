package grpc

import (
	"context"
	"profile/internal/service"
	"profile/pkg/profilepb"
)

type ProfileHandler struct {
	profilepb.UnimplementedProfileServiceServer
	s service.ProfileService
}

func NewProfileHandler(s service.ProfileService) *ProfileHandler {
	return &ProfileHandler{s: s}
}

func (h *ProfileHandler) GetAll(ctx context.Context, req *profilepb.GetAllRequest) (*profilepb.GetAllResponse, error) {
	profiles := h.s.GetAll(ctx)
	resp := &profilepb.GetAllResponse{Profiles: nil}
	for _, val := range profiles {
		dto := &profilepb.ProfileDTO{
			UserId:      val.UserId,
			Username:    "test",
			Name:        val.Name,
			Gender:      val.Gender,
			Description: val.Description,
			PhotoPath:   val.PhotoPath,
		}
		resp.Profiles = append(resp.Profiles, dto)
	}
	return resp, nil
}
