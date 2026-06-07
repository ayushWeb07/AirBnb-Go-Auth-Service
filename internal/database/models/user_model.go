package models

type UserModel struct {
	ID       string `json:"id"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
}
