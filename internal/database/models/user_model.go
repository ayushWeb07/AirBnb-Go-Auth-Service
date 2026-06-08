package models

type UserModel struct {
	ID        string `db:"id"`
	Username  string `db:"username" json:"username" validate:"required"`
	Email     string `db:"email" json:"email" validate:"required"`
	Password  string `db:"password" json:"password" validate:"required"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}
