package dtos

type GetPermissionsOfUserPayload struct {
	UserID int `json:"user_id" validate:"required,gte=1"`
}

type CheckUserHasPermissionPayload struct {
	UserID       int `json:"user_id" validate:"required,gte=1"`
	PermissionID int `json:"permission_id" validate:"required,gte=1"`
}

type GetPermissionsOfRolePayload struct {
	RoleID int `json:"role_id" validate:"required,gte=1"`
}

type AssignPermissionToRolePayload struct {
	RoleID       int `json:"role_id" validate:"required,gte=1"`
	PermissionID int `json:"permission_id" validate:"required,gte=1"`
}

type RemovePermissionFromRolePayload struct {
	RoleID       int `json:"role_id" validate:"required,gte=1"`
	PermissionID int `json:"permission_id" validate:"required,gte=1"`
}

type CheckRoleHasPermissionPayload struct {
	RoleID       int `json:"role_id" validate:"required,gte=1"`
	PermissionID int `json:"permission_id" validate:"required,gte=1"`
}
