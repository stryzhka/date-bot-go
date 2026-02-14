package http

import (
	"bytes"
	"date-bot-go/profile"
	"date-bot-go/profile/models"
	"date-bot-go/profile/repository/mock"
	"date-bot-go/profile/service"
	profileMock "date-bot-go/profile/service/mock"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
)

func TestHealthcheck(t *testing.T) {
	r := new(mock.MockRepository)
	s := service.NewProfileService(r)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	w := httptest.NewRecorder()
	h.Healthcheck(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestEmptyGetAll(t *testing.T) {
	s := new(profileMock.MockService)
	s.On("GetAll").Return([]models.Profile{})
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodGet, "/api/profile/", nil)
	w := httptest.NewRecorder()
	h.GetAll(w, req)
	assert.Equal(t, "[]", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSuccessGetAll(t *testing.T) {
	s := new(profileMock.MockService)
	testProfile := &models.Profile{
		Id:          uuid.Nil,
		UserId:      "123",
		Name:        "test",
		Gender:      "f",
		Description: "test",
		Topics:      nil,
		DateCreated: time.Time{},
		PhotoPath:   "test",
	}
	s.On("GetAll").Return([]models.Profile{*testProfile, *testProfile, *testProfile})
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodGet, "/api/profile/", nil)
	w := httptest.NewRecorder()
	h.GetAll(w, req)
	//assert.Equal(t, "[]", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	log.Println(w.Body.String())
}

func TestEmptyGetById(t *testing.T) {
	s := new(profileMock.MockService)
	s.On("GetById", "0000-0000-0000-0000").Return(nil)
	h := NewHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/profile/{id}", h.GetById)

	req := httptest.NewRequest(http.MethodGet, "/api/profile/0000-0000-0000-0000", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "null", w.Body.String())
	s.AssertExpectations(t)
}

func TestSuccessGetById(t *testing.T) {
	profile := &models.Profile{
		Id:          uuid.Nil,
		UserId:      "123",
		Name:        "test",
		Gender:      "f",
		Description: "test test test",
		Topics:      nil,
		DateCreated: time.Time{},
		PhotoPath:   "test",
	}
	s := new(profileMock.MockService)
	s.On("GetById", "0000-0000-0000-0000").Return(profile)
	h := NewHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/profile/{id}", h.GetById)

	expJsonProfile, err := json.Marshal(profile)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodGet, "/api/profile/0000-0000-0000-0000", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, bytes.NewBuffer(expJsonProfile).String(), w.Body.String())
	s.AssertExpectations(t)
}

func TestFailValidationCreate(t *testing.T) {
	//incorrectProfileDto := &ProfileDto{
	//	UserId:      "",
	//	Name:        "",
	//	Gender:      "",
	//	Description: "",
	//	PhotoPath:   "",
	//}
	jsonIncorrectProfileDto := "{iuojo"
	log.Println(jsonIncorrectProfileDto)
	//assert.NoError(t, err)
	s := new(profileMock.MockService)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodPost, "/api/profile/", bytes.NewBufferString(jsonIncorrectProfileDto))
	w := httptest.NewRecorder()
	s.On("Create", "", "", "", "", "").Return(errors.New("test error"))
	h.Create(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	log.Println(w.Body.String())
}

func TestFailValidationGenderCreate(t *testing.T) {
	incorrectProfileDto := &ProfileDto{
		UserId:      "",
		Name:        "",
		Gender:      "",
		Description: "",
		PhotoPath:   "",
	}
	jsonIncorrectProfileDto, err := json.Marshal(incorrectProfileDto)
	log.Println(string(jsonIncorrectProfileDto))
	assert.NoError(t, err)
	s := new(profileMock.MockService)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodPost, "/api/profile/", bytes.NewReader(jsonIncorrectProfileDto))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.On("Create", "", "", "", "", "").Return(errors.New("test error"))
	h.Create(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	log.Println(w.Body.String())
}

func TestFailValidationNameCreate(t *testing.T) {
	incorrectProfileDto := &ProfileDto{
		UserId:      "",
		Name:        "",
		Gender:      "f",
		Description: "",
		PhotoPath:   "",
	}
	jsonIncorrectProfileDto, err := json.Marshal(incorrectProfileDto)
	log.Println(string(jsonIncorrectProfileDto))
	assert.NoError(t, err)
	s := new(profileMock.MockService)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodPost, "/api/profile/", bytes.NewReader(jsonIncorrectProfileDto))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.On("Create", "", "", "", "", "").Return(errors.New("test error"))
	h.Create(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	log.Println(w.Body.String())
}

func TestFailValidationUserIdCreate(t *testing.T) {
	incorrectProfileDto := &ProfileDto{
		UserId:      "",
		Name:        "test",
		Gender:      "f",
		Description: "",
		PhotoPath:   "",
	}
	jsonIncorrectProfileDto, err := json.Marshal(incorrectProfileDto)
	log.Println(string(jsonIncorrectProfileDto))
	assert.NoError(t, err)
	s := new(profileMock.MockService)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodPost, "/api/profile/", bytes.NewReader(jsonIncorrectProfileDto))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.On("Create", "", "", "", "", "").Return(errors.New("test error"))
	h.Create(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	log.Println(w.Body.String())
}

func TestFailLogicCreate(t *testing.T) {
	incorrectProfileDto := &ProfileDto{
		UserId:      "123",
		Name:        "test",
		Gender:      "f",
		Description: "",
		PhotoPath:   "",
	}
	jsonIncorrectProfileDto, err := json.Marshal(incorrectProfileDto)
	log.Println(string(jsonIncorrectProfileDto))
	assert.NoError(t, err)
	s := new(profileMock.MockService)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodPost, "/api/profile/", bytes.NewReader(jsonIncorrectProfileDto))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.On("Create", mock2.Anything, "123", "test", "f", "").Return(profile.ErrUserAlreadyExists)
	h.Create(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	log.Println(w.Body.String())
}

func TestSuccessCreate(t *testing.T) {
	incorrectProfileDto := &ProfileDto{
		UserId:      "123",
		Name:        "test",
		Gender:      "f",
		Description: "",
		PhotoPath:   "",
	}
	jsonIncorrectProfileDto, err := json.Marshal(incorrectProfileDto)
	log.Println(string(jsonIncorrectProfileDto))
	assert.NoError(t, err)
	s := new(profileMock.MockService)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodPost, "/api/profile/", bytes.NewReader(jsonIncorrectProfileDto))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.On("Create", mock2.Anything, "123", "test", "f", "").Return(nil)
	h.Create(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	log.Println(w.Body.String())
}

func TestFailUserNotFoundUpdate(t *testing.T) {
	incorrectProfileDto := &ProfileDto{
		UserId:      "123",
		Name:        "test",
		Gender:      "f",
		Description: "",
		PhotoPath:   "",
	}
	expProfile := &models.Profile{
		Id:          uuid.Nil,
		UserId:      "123",
		Name:        "test",
		Gender:      "f",
		Description: "",
		Topics:      nil,
		DateCreated: time.Time{},
		PhotoPath:   "",
	}
	jsonIncorrectProfileDto, err := json.Marshal(incorrectProfileDto)
	log.Println(string(jsonIncorrectProfileDto))
	assert.NoError(t, err)
	s := new(profileMock.MockService)
	s.On("UpdateById", mock2.Anything, "", expProfile).Return(profile.ErrUserNotFound)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodPost, "/api/profile/", bytes.NewReader(jsonIncorrectProfileDto))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Update(w, req)
	log.Println(w.Body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFailValidationUpdate(t *testing.T) {
	//incorrectProfileDto := &ProfileDto{
	//	UserId:      "",
	//	Name:        "",
	//	Gender:      "",
	//	Description: "",
	//	PhotoPath:   "",
	//}
	incorrectProfileDto := &ProfileDto{
		UserId:      "123",
		Name:        "test",
		Gender:      "f",
		Description: "",
		PhotoPath:   "",
	}
	expProfile := &models.Profile{
		Id:          uuid.Nil,
		UserId:      "123",
		Name:        "test",
		Gender:      "f",
		Description: "",
		Topics:      nil,
		DateCreated: time.Time{},
		PhotoPath:   "",
	}
	jsonIncorrectProfileDto := "{iuojo"
	log.Println(jsonIncorrectProfileDto)
	//assert.NoError(t, err)
	s := new(profileMock.MockService)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodPost, "/api/profile/", bytes.NewBufferString(jsonIncorrectProfileDto))
	w := httptest.NewRecorder()
	s.On("Create", "", "", "", "", "").Return(errors.New("test error"))
	h.Create(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	log.Println(w.Body.String())
}

func TestFailValidationGenderCreate(t *testing.T) {
	incorrectProfileDto := &ProfileDto{
		UserId:      "",
		Name:        "",
		Gender:      "",
		Description: "",
		PhotoPath:   "",
	}
	jsonIncorrectProfileDto, err := json.Marshal(incorrectProfileDto)
	log.Println(string(jsonIncorrectProfileDto))
	assert.NoError(t, err)
	s := new(profileMock.MockService)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodPost, "/api/profile/", bytes.NewReader(jsonIncorrectProfileDto))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.On("Create", "", "", "", "", "").Return(errors.New("test error"))
	h.Create(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	log.Println(w.Body.String())
}

func TestFailValidationNameCreate(t *testing.T) {
	incorrectProfileDto := &ProfileDto{
		UserId:      "",
		Name:        "",
		Gender:      "f",
		Description: "",
		PhotoPath:   "",
	}
	jsonIncorrectProfileDto, err := json.Marshal(incorrectProfileDto)
	log.Println(string(jsonIncorrectProfileDto))
	assert.NoError(t, err)
	s := new(profileMock.MockService)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodPost, "/api/profile/", bytes.NewReader(jsonIncorrectProfileDto))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.On("Create", "", "", "", "", "").Return(errors.New("test error"))
	h.Create(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	log.Println(w.Body.String())
}

func TestFailValidationUserIdCreate(t *testing.T) {
	incorrectProfileDto := &ProfileDto{
		UserId:      "",
		Name:        "test",
		Gender:      "f",
		Description: "",
		PhotoPath:   "",
	}
	jsonIncorrectProfileDto, err := json.Marshal(incorrectProfileDto)
	log.Println(string(jsonIncorrectProfileDto))
	assert.NoError(t, err)
	s := new(profileMock.MockService)
	h := NewHandler(s)
	req := httptest.NewRequest(http.MethodPost, "/api/profile/", bytes.NewReader(jsonIncorrectProfileDto))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.On("Create", "", "", "", "", "").Return(errors.New("test error"))
	h.Create(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	log.Println(w.Body.String())
}
