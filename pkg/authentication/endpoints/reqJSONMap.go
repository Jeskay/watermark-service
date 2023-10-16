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

type GenerateRequest struct{}

type VerifyTwoFactorRequest struct {
	Token string `json:"token"`
}

type ValidateRequest struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type VerifyJwtRequest struct {
	Token string `json:"token"`
}

type DisableRequest struct{}
type RegisterResponse struct {
	UserId string `json:"userId"`
	Err    string `json:"err,omitempty"`
}

type LoginResponse struct {
	Status int64  `json:"status"`
	Token  string `json:"token"`
}

type GenerateResponse struct {
	Base32     string `json:"base32"`
	OtpAuthUrl string `json:"otp_auth_url"`
}

type VerifyTwoFactorResponse struct {
	OtpVerified bool           `json:"otp_verified"`
	User        *internal.User `json:"user"`
}

type ValidateResponse struct {
	OtpValid bool `json:"otp_valid"`
}

type VerifyJwtResponse struct {
	Verified bool           `json:"verified"`
	User     *internal.User `json:"user"`
}

type DisableResponse struct {
	OtpDisabled bool           `json:"otp_disabled"`
	User        *internal.User `json:"user"`
}

type ServiceStatusRequest struct{}

type ServiceStatusResponse struct {
	Code int    `json:"code"`
	Err  string `json:"err,omitempty"`
}
