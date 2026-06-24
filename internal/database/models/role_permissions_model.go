package models

type RolePermissionModel struct {
	ID           int    `db:"id" json:"id"`
	RoleID       int    `db:"role_id" json:"role_id"`
	PermissionID int    `db:"permission_id" json:"permission_id"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	UpdatedAt    string `db:"updated_at" json:"updated_at"`
}
