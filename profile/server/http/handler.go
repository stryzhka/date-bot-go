package http

import (
	"date-bot-go/profile"
	"date-bot-go/profile/models"
	"encoding/json"
	"net/http"
)

type Handler struct {
	s profile.Service
}

func NewHandler(s profile.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Healthcheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// GetAll godoc
// @Summary Get all profiles
// @Tags profile
// // @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.Profile
// // @Failure 401
// @Router /api/profile/ [get]
func (h *Handler) GetAll(w http.ResponseWriter, req *http.Request) {
	var profiles []models.Profile
	profiles = h.s.GetAll(req.Context())
	jsonProfiles, err := json.Marshal(profiles)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		jsonProfiles = nil
	}
	w.Write(jsonProfiles)
}

//func (h *Handler) Create(w http.ResponseWriter, req *http.Request) {
//	ctx := req.Context()
//}
