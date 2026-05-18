package models

import (
	"github.com/google/uuid"
	"time"
)

type Profile struct {
	Id          uuid.UUID `json:"id"`
	UserId      string    `json:"user_id"`
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	Description string    `json:"description"`
	Topics      []Topic   `json:"topics"`
	DateCreated time.Time `json:"date_created"`
	PhotoPath   string    `json:"photo_path"`
}
