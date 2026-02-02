package models

type Profile struct {
	//Id          uuid.UUID `json:"id"`
	UserId      string `json:"user_id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Description string `json:"description"`
	//Topics      []Topic   `json:"topics"`
	//DateCreated time.Time `json:"dateCreated"`
	PhotoPath string `json:"photoPath"`
}
