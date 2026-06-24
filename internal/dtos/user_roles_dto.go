package dtos

type GetRolesOfUserPayload struct {
	UserID int `json:"user_id" validate:"required,gte=1"`
}

type AssignRoleToUserPayload struct {
	UserID int `json:"user_id" validate:"required,gte=1"`
	RoleID int `json:"role_id" validate:"required,gte=1"`
}

type RemoveRoleFromUserPayload struct {
	UserID int `json:"user_id" validate:"required,gte=1"`
	RoleID int `json:"role_id" validate:"required,gte=1"`
}

type CheckUserHasRolePayload struct {
	UserID int `json:"user_id" validate:"required,gte=1"`
	RoleID int `json:"role_id" validate:"required,gte=1"`
}

type CheckUserHasAllRolesPayload struct {
	UserID    int      `json:"user_id" validate:"required,gte=1"`
	RoleNames []string `json:"role_names" validate:"required,min=1,unique"`
}

type CheckUserHasAnyRolesPayload struct {
	UserID    int      `json:"user_id" validate:"required,gte=1"`
	RoleNames []string `json:"role_names" validate:"required,min=1,unique"`
}
