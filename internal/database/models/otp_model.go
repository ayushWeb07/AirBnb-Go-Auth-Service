package models

type OtpModel struct {
	ID        int    `db:"id" json:"id"`
	UserID    int    `db:"user_id" json:"user_id" validate:"required"`
	UserEmail string `db:"user_email" json:"user_email" validate:"required"`
	OtpHash   string `db:"otp_hash" json:"otp_hash" validate:"required"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}
