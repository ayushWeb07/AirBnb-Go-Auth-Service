package models

type RoleModel struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name" validate:"required"`
	Description string `db:"description" json:"description" validate:"required"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}
