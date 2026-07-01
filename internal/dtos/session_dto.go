package dtos

type CreateSessionPayload struct {
	UserID           int    `json:"user_id" validate:"required,gte=1"`
	RefreshTokenHash string `json:"refresh_token_hash" validate:"required,min=6"`
}

type FetchSessionPayload struct {
	UserID           int    `json:"user_id" validate:"required,gte=1"`
	RefreshTokenHash string `json:"refresh_token_hash" validate:"required,min=6"`
	Revoked          bool   `json:"revoked"`
}

type UpdateSessionByIdParams struct {
	ID int `json:"id" validate:"required,gte=1"`
}

type UpdateSessionByIdPayload struct {
	Revoked bool `json:"revoked" validate:"required"`
}

type RevokeAllSessionsPayload struct {
	UserID int `json:"user_id" validate:"required,gte=1"`
}
