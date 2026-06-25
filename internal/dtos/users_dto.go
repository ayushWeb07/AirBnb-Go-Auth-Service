package dtos

type CreateUserPayload struct {
	Username string `json:"username" validate:"required,min=6,max=100"`
	Email    string `json:"email" validate:"required,email,min=6,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
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
	Token string `json:"token" validate:"required"`
}

type LogoutUserPayload struct {
	Token string `json:"token" validate:"required"`
}

type LogoutUserFromAllSessionsPayload struct {
	Token string `json:"token" validate:"required"`
}

type SendOtpForVerificationPayload struct {
	UserEmail string `json:"user_email" validate:"required"`
}

type VerifyOtpPayload struct {
	UserEmail string `json:"user_email" validate:"required"`
	Otp       string `json:"otp" validate:"required"`
}
