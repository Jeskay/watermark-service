package endpoints

import "watermark-service/internal"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GenerateRequest struct {
	UserId string `json:"user_id"`
}

type RegisterResponse struct {
	UserId string `json:"userId"`
	Err    string `json:"err,omitempty"`
}

type LoginResponse struct {
	Status int64          `json:"status"`
	User   *internal.User `json:"user"`
}

type GenerateResponse struct {
	Base32     string `json:"base32"`
	OtpAuthUrl string `json:"otp_auth_url"`
}

type ServiceStatusRequest struct{}

type ServiceStatusResponse struct {
	Code int    `json:"code"`
	Err  string `json:"err,omitempty"`
}
