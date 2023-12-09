package internal

type User struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	OtpEnabled bool   `json:"otp_enabled"`
}
