package models

type Like struct {
	UserId  string `json:"userId"`
	LikedId string `json:"likedId"`
}
