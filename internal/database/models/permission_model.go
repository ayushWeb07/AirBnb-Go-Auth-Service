package models

type PermissionModel struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name" validate:"required"`
	Description string `db:"description" json:"description" validate:"required"`
	Resource    string `db:"resource" json:"resource" validate:"required"`
	Action      string `db:"action" json:"action" validate:"required"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}
