package models

type UserModel struct {
	ID        int    `db:"id" json:"id"`
	Username  string `db:"username" json:"username" validate:"required"`
	Email     string `db:"email" json:"email" validate:"required"`
	Password  string `db:"password" json:"-" validate:"required"`
	Verified  bool   `db:"verified" json:"verified"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	UpdatedAt string `db:"updated_at" json:"updatedAt"`
}
