package http

import (
	"date-bot-go/profile"
	"date-bot-go/profile/models"
	"encoding/json"
	"net/http"
	"strings"
)

type errorMessage struct {
	Error string `json:"error"`
}

type message struct {
	Message string `json:"message"`
}

type ProfileDto struct {
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Description string `json:"description"`
	PhotoPath   string `json:"photo_path"`
}

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
	return
}

// GetById godoc
// @Summary Get profile by id
// @Tags profile
// // @Security ApiKeyAuth
// @Produce json
// @Success 200 models.Profile
// // @Failure 401
// @Router /api/profile/{id} [get]
func (h *Handler) GetById(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	jsonProfile, err := json.Marshal(h.s.GetById(req.Context(), id))
	if err != nil {
		//TODO: это неправильно
		jsonProfile = nil
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonProfile)
}

// Create godoc
// @Summary Create new profile
// @Tags profile
// // @Security ApiKeyAuth
// @Produce json
// @Success 200 models.Profile
// // @Failure 401
// @Failure 400
// @Router /api/profile/ [post]
func (h *Handler) Create(w http.ResponseWriter, req *http.Request) {
	//TODO: что то здесь не так...
	var parsedProfileDto ProfileDto
	err := json.NewDecoder(req.Body).Decode(&parsedProfileDto)
	if err != nil {
		msg, _ := json.Marshal(&errorMessage{Error: "can't parse json"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	if len(parsedProfileDto.Name) >= 15 {
		msg, _ := json.Marshal(&errorMessage{Error: "name is too long"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	if strings.TrimSpace(parsedProfileDto.Gender) == "" {
		msg, _ := json.Marshal(&errorMessage{Error: "gender can't be empty"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	if strings.TrimSpace(parsedProfileDto.Name) == "" {
		msg, _ := json.Marshal(&errorMessage{Error: "name can't be empty"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}

	if strings.TrimSpace(parsedProfileDto.UserId) == "" {
		msg, _ := json.Marshal(&errorMessage{Error: "user_id can't be empty"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	if len(parsedProfileDto.Description) >= 128 {
		msg, _ := json.Marshal(&errorMessage{Error: "description is too long"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	err = h.s.Create(req.Context(), parsedProfileDto.UserId, parsedProfileDto.Name, parsedProfileDto.Gender, parsedProfileDto.Description)
	if err != nil {
		msg, _ := json.Marshal(&errorMessage{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	w.WriteHeader(http.StatusCreated) //TODO: но ваще кстати можно просто отправить созданный профиль
}

//func (h *Handler) Create(w http.ResponseWriter, req *http.Request) {
//	ctx := req.Context()
//}
