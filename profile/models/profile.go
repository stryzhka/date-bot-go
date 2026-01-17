package models

import "github.com/google/uuid"

type Profile struct {
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	Description string    `json:"description"`
	Topics      []Topic   `json:"topics"`
}
