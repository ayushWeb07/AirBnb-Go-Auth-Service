package models

type UserRoleModel struct {
	ID        int    `db:"id" json:"id"`
	RoleID    int    `db:"role_id" json:"role_id"`
	UserID    int    `db:"user_id" json:"user_id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}
