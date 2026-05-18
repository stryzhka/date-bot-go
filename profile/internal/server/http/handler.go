package http

import (
	"encoding/json"
	"net/http"
	profile "profile/internal"
	"profile/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func validate(parsedProfileDto *ProfileDto) error {
	if len(parsedProfileDto.Name) >= 15 {
		return profile.ErrValidationName
	}
	if strings.TrimSpace(parsedProfileDto.Gender) == "" {
		return profile.ErrValidationGender
	}
	if strings.TrimSpace(parsedProfileDto.Name) == "" {
		return profile.ErrValidationName
	}

	if strings.TrimSpace(parsedProfileDto.UserId) == "" {
		return profile.ErrValidationUserId
	}
	if len(parsedProfileDto.Description) >= 128 {
		return profile.ErrValidationDescription
	}
	return nil
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err != nil {
		jsonProfiles = []byte("[]")
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
	id := mux.Vars(req)["id"]
	jsonProfile, err := json.Marshal(h.s.GetById(req.Context(), id))
	if err != nil {
		jsonProfile = []byte("{}")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonProfile)
}

// Create godoc
// @Summary Create new profile
// @Tags profile
// // @Security ApiKeyAuth
// @Accepts json ProfileDto
// @Produce json
// @Success 201
// // @Failure 401
// @Failure 400
// @Router /api/profile/ [post]
func (h *Handler) Create(w http.ResponseWriter, req *http.Request) {

	//TODO: вроде ща все так...
	var parsedProfileDto ProfileDto
	err := json.NewDecoder(req.Body).Decode(&parsedProfileDto)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		msg, _ := json.Marshal(&errorMessage{Error: "can't parse json"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	err = validate(&parsedProfileDto)
	if err != nil {
		msg, _ := json.Marshal(&errorMessage{Error: err.Error()})
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

// Update godoc
// @Summary Update profile by id
// @Tags profile
// // @Security ApiKeyAuth
// @Accepts json ProfileDto
// @Produce json
// @Success 200
// // @Failure 401
// @Failure 400
// @Router /api/profile/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	var parsedProfileDto ProfileDto
	err := json.NewDecoder(req.Body).Decode(&parsedProfileDto)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		msg, _ := json.Marshal(&errorMessage{Error: "can't parse json"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	err = validate(&parsedProfileDto)
	if err != nil {
		msg, _ := json.Marshal(&errorMessage{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	//todo: кажется тут вообще надо просто параметры а не модель
	newProfile := &models.Profile{
		Id:          uuid.Nil,
		UserId:      parsedProfileDto.UserId,
		Name:        parsedProfileDto.Name,
		Gender:      parsedProfileDto.Gender,
		Description: parsedProfileDto.Description,
		Topics:      nil,
		DateCreated: time.Time{},
		PhotoPath:   parsedProfileDto.PhotoPath,
	}
	err = h.s.UpdateById(req.Context(), id, newProfile)
	if err != nil {
		msg, _ := json.Marshal(&errorMessage{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Delete godoc
// @Summary Delete profile by id
// @Tags profile
// // @Security ApiKeyAuth
// @Accepts json ProfileDto
// @Produce json
// @Success 200
// // @Failure 401
// @Failure 400
// @Router /api/profile/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	err := h.s.DeleteById(req.Context(), id)
	var msg []byte
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		msg, _ = json.Marshal(errorMessage{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}
	w.WriteHeader(http.StatusOK)
}
