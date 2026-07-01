package dtos

type CreateUserPayload struct {
	Username string `json:"username" validate:"required,min=6,max=100"`
	Email    string `json:"email" validate:"required,email,min=6,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UpdateUserByIdParams struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type UpdateUserByIdPayload struct {
	Verified bool `json:"verified" validate:"required"`
}

type GetUserByEmailPayload struct {
	Email string `json:"email" validate:"required,email,min=6,max=100"`
}

type GetUserByUsernameAndEmailPayload struct {
	Username string `json:"username" validate:"required,min=6,max=100"`
	Email    string `json:"email" validate:"required,email,min=6,max=100"`
}

type LoginUserPayload struct {
	Username string `json:"username" validate:"required,min=6,max=100"`
	Email    string `json:"email" validate:"required,email,min=6,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type GetUserByIdParams struct {
	ID int `validate:"required,gte=1"`
}

type DeleteUserByIdParams struct {
	ID int `validate:"required,gte=1"`
}

type RefreshAccessTokenPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutUserPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutUserFromAllSessionsPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
