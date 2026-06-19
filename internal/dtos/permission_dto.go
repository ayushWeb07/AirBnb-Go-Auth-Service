package dtos

type CreatePermissionPayload struct {
	Name        string `json:"name" validate:"required,min=6,max=100"`
	Description string `json:"description" validate:"required,min=6,max=100"`
	Resource    string `json:"resource" validate:"required,min=3,max=50"`
	Action      string `json:"action" validate:"required,min=3,max=50"`
}

type GetPermissionByIdParams struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type DeletePermissionByIdParams struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type UpdatePermissionByIdParams struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type UpdatePermissionByIdPayload struct {
	Name        string `json:"name" validate:"min=6,max=100"`
	Description string `json:"description" validate:"min=6,max=100"`
	Resource    string `json:"resource" validate:"min=3,max=50"`
	Action      string `json:"action" validate:"min=3,max=50"`
}
