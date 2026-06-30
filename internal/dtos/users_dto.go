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

type LogoutUserFromAllSessionsPayload struct {
	Token string `json:"token" validate:"required"`
}

type SendOtpForVerificationPayload struct {
	UserEmail string `json:"user_email" validate:"required"`
}

type CreateSessionPayload struct {
	UserID           int    `json:"user_id" validate:"required,gte=1"`
	RefreshTokenHash string `json:"refresh_token_hash" validate:"required,min=6"`
}

type CreateOtpServicePayload struct {
	UserEmail string `json:"user_email" validate:"required,email,min=6,max=100"`
}

type CreateOtpRepoPayload struct {
	UserID    int    `json:"user_id" validate:"required,gte=1"`
	UserEmail string `json:"user_email" validate:"required,email,min=6,max=100"`
	OtpHash   string `json:"otp_hash" validate:"required,min=6"`
}

type VerifyOtpPayload struct {
	UserEmail string `json:"user_email" validate:"required,email,min=6,max=100"`
	Otp       string `json:"otp" validate:"required,len=10"`
}

type FetchOtpRepoPayload struct {
	UserID    int    `json:"user_id" validate:"required,gte=1"`
	UserEmail string `json:"user_email" validate:"required,email,min=6,max=100"`
	OtpHash   string `json:"otp_hash" validate:"required,min=6"`
}

type DeleteOtpsRepoPayload struct {
	UserID int `json:"user_id" validate:"required,gte=1"`
}

type RefreshAccessTokenPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type FetchSessionPayload struct {
	UserID           int    `json:"user_id" validate:"required,gte=1"`
	RefreshTokenHash string `json:"refresh_token_hash" validate:"required,min=6"`
	Revoked          bool   `json:"revoked"`
}

type LogoutUserPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdateSessionByIdParams struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type UpdateSessionByIdPayload struct {
	Revoked bool `json:"revoked" validate:"required"`
}
