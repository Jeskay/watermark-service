package internal

import (
	uuid "github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	OtpEnabled bool      `json:"otp_enabled"`
}
