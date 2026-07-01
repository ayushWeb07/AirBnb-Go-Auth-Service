package dtos

type SendOtpForVerificationPayload struct {
	UserEmail string `json:"user_email" validate:"required"`
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
