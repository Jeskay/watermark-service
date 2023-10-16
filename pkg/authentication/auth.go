package authentication

import (
	"context"
	"errors"
	"net/http"
	"time"

	"strings"
	"watermark-service/internal"
	auth "watermark-service/internal/authentication"

	jwt "github.com/golang-jwt/jwt/v4"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type authService struct {
	orm        *gorm.DB
	signingKey []byte
}

type userClaims struct {
	internal.User
	jwt.RegisteredClaims
}

func NewService(dbORM *gorm.DB, signingKey string) *authService {
	return &authService{orm: dbORM, signingKey: []byte(signingKey)}
}

func (a *authService) Login(ctx context.Context, email, password string) (int64, string) {
	var user auth.User
	token2FA := ctx.Value("2FA").(string)
	result := a.orm.First(&user, "email = ?", email)
	if result.Error != nil {
		return http.StatusUnauthorized, ""
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return http.StatusUnauthorized, ""
	}
	if user.Otp_enabled && user.Otp_verified && !totp.Validate(token2FA, user.Otp_secret) {
		return http.StatusUnauthorized, ""
	}
	claims := userClaims{
		internal.User{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			Issuer:    "watermark-service",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.signingKey)
	if err != nil {
		return http.StatusUnauthorized, ""
	}
	return http.StatusAccepted, signedToken
}

func (a *authService) Register(_ context.Context, email, name, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
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

func (a *authService) Generate(ctx context.Context) (string, string) {
	var user auth.User
	claimedUser, ok := ctx.Value("user").(*internal.User)
	if !ok {
		return "", ""
	}
	result := a.orm.First(&user, "id = ?", claimedUser.ID)
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
		Otp_enabled:  true,
		Otp_secret:   key.Secret(),
		Otp_auth_url: key.URL(),
	}

	a.orm.Model(&user).Updates(updateData)
	return key.Secret(), key.URL()
}

func (a *authService) Verify(ctx context.Context, token string) (bool, *internal.User) {
	var user auth.User
	claimedUser, ok := ctx.Value("user").(*internal.User)
	if !ok {
		return false, nil
	}
	result := a.orm.First(&user, "id = ?", claimedUser.ID)
	if result.Error != nil {
		return false, nil
	}
	isValid := totp.Validate(token, user.Otp_secret)
	if !isValid {
		return false, nil
	}
	if user.Otp_enabled && !user.Otp_verified {
		updateData := auth.User{
			Otp_verified: true,
		}
		a.orm.Model(&user).Updates(updateData)
	}
	userResp := internal.User{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		OtpEnabled: user.Otp_enabled,
	}
	return true, &userResp
}

func (a *authService) VerifyJwt(_ context.Context, tokenString string) (bool, *internal.User) {
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.signingKey, nil
	})

	if err != nil {
		return false, nil
	}

	if claims, ok := token.Claims.(*userClaims); ok && token.Valid {
		return true, &internal.User{
			ID:         claims.User.ID,
			Name:       claims.Name,
			Email:      claims.Email,
			OtpEnabled: claims.OtpEnabled,
		}
	}
	return false, nil
}

func (a *authService) Validate(_ context.Context, userID, token string) bool {
	var user auth.User
	var unmarshalled uuid.UUID
	err := unmarshalled.UnmarshalBinary([]byte(userID))
	if err != nil {
		return false
	}
	result := a.orm.First(&user, "id = ?", unmarshalled)
	if result.Error != nil {
		return false
	}
	return totp.Validate(token, user.Otp_secret)
}

func (a *authService) Disable(ctx context.Context) (bool, *internal.User) {
	var user auth.User
	claimedUser, ok := ctx.Value("user").(internal.User)
	if !ok {
		return false, nil
	}
	result := a.orm.First(&user, "id = ?", claimedUser.ID)
	if result.Error != nil {
		return false, nil
	}
	user.Otp_enabled = false
	a.orm.Save(&user)
	userResp := internal.User{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		OtpEnabled: user.Otp_enabled,
	}
	return user.Otp_enabled, &userResp
}

func (a *authService) ServiceStatus(_ context.Context) (int, error) {
	return http.StatusOK, nil
}
