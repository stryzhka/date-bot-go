package http

import (
	"date-bot-go/profile/models"
	"date-bot-go/profile/repository/mock"
	"date-bot-go/profile/service"
	profileMock "date-bot-go/profile/service/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
