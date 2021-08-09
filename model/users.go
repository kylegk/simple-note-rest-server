package model

type UserAccount struct {
	UserID int
	User string
}

type CreateUserRequest struct {
	User string `json:"user"`
}

type CreateUserResponse struct {
	UserID int `json:"userid"`
	Token string `json:"token"`
}

