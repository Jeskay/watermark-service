package authentication

import (
	"context"
	"errors"
	"net/http"

	"strings"
	"watermark-service/internal"
	auth "watermark-service/internal/authentication"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type authService struct {
	orm *gorm.DB
}

func NewService(dbORM *gorm.DB) *authService {
	return &authService{orm: dbORM}
}

func (a *authService) Login(_ context.Context, email, password string) (int64, *internal.User) {
	var user auth.User
	result := a.orm.First(&user, "email = ?", email)
	if result.Error != nil {
		return http.StatusUnauthorized, nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return http.StatusUnauthorized, nil
	}
	userResp := &internal.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	return http.StatusAccepted, userResp
}

func (a *authService) Register(_ context.Context, email, name, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		return "", err
	}
	newUser := auth.User{
		Name:     name,
		Email:    email,
		Password: string(hash),
	}
	result := a.orm.Create(&newUser)
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return "", errors.New("User already exists")
	} else if result.Error != nil {
		return "", errors.New(result.Error.Error())
	}
	return newUser.ID.String(), nil
}

func (a *authService) Generate(_ context.Context, userId string) (string, string) {
	var user auth.User
	var unmarshalled uuid.UUID
	err := unmarshalled.UnmarshalBinary([]byte(userId))
	if err != nil {
		return "", ""
	}
	result := a.orm.First(&user, "id = ?", unmarshalled)
	if result.Error != nil {
		return "", ""
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "watermark-service",
		AccountName: user.Email,
		SecretSize:  16,
	})
	if err != nil {
		return "", ""
	}
	updateData := auth.User{
		Otp_secret:   key.Secret(),
		Otp_auth_url: key.URL(),
	}

	a.orm.Model(&user).Updates(updateData)
	return key.Secret(), key.URL()
}

func (a *authService) ServiceStatus(_ context.Context) (int, error) {
	return http.StatusOK, nil
}
