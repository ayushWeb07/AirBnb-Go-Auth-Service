package dtos

type CreateRolePayload struct {
	Name        string `json:"name" validate:"required,min=6,max=100"`
	Description string `json:"description" validate:"required,min=6,max=100"`
}

type GetRoleByIdParams struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type DeleteRoleByIdParams struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type UpdateRoleByIdParams struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type UpdateRoleByIdPayload struct {
	Name        string `json:"name" validate:"required,min=6,max=100"`
	Description string `json:"description" validate:"required,min=6,max=100"`
}
