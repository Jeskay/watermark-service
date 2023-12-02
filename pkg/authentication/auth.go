package authentication

import (
	"context"
	"errors"
	"net/http"
	"time"

	"strings"
	"watermark-service/internal"
	"watermark-service/internal/authentication"
	auth "watermark-service/internal/authentication"

	jwt "github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type authService struct {
	ORMInstance *gorm.DB
	DBAvailable bool
	Dsn         string
	SigningKey  []byte
	log         *zap.Logger
}

type userClaims struct {
	internal.User
	jwt.RegisteredClaims
}

func NewService(dbConnection internal.DatabaseConnectionStr, signingKey string) *authService {
	dsn := dbConnection.GetDSN()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}))
	service := &authService{
		ORMInstance: db,
		DBAvailable: true,
		SigningKey:  []byte(signingKey),
		Dsn:         dsn,
		log:         zap.L().With(zap.String("Service", "AuthenticationService")),
	}
	if err == nil {
		err = authentication.InitDb(db)
	}
	if err != nil {
		service.DBAvailable = false
		go service.Reconnect()
	}
	return service
}

func (a *authService) Reconnect() {
	for idleTime := time.Duration(2); idleTime < 1000; idleTime *= idleTime {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN: a.Dsn,
		}))
		if err == nil {
			err = authentication.InitDb(db)
		}
		if err == nil {
			a.ORMInstance = db
			a.DBAvailable = true
			break
		}
		a.log.Error("Reconnect", zap.String("Database", "Failed"), zap.Error(err))
		a.log.Info("Reconnect", zap.String("Database", "Attempt"), zap.Duration("After", idleTime))
		time.Sleep(idleTime * time.Second)
	}
	a.log.Info("Reconnect", zap.String("Status", "Success"), zap.String("Connection", a.Dsn))
}

func (a *authService) Login(ctx context.Context, email, password string) (int64, string) {
	var user auth.User
	token2FA := ctx.Value("2FA").(string)
	result := a.ORMInstance.First(&user, "email = ?", email)
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
	signedToken, err := token.SignedString(a.SigningKey)
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
	result := a.ORMInstance.Create(&newUser)
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
	result := a.ORMInstance.First(&user, "id = ?", claimedUser.ID)
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

	a.ORMInstance.Model(&user).Updates(updateData)
	return key.Secret(), key.URL()
}

func (a *authService) Verify(ctx context.Context, token string) (bool, *internal.User) {
	var user auth.User
	claimedUser, ok := ctx.Value("user").(*internal.User)
	if !ok {
		return false, nil
	}
	result := a.ORMInstance.First(&user, "id = ?", claimedUser.ID)
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
		a.ORMInstance.Model(&user).Updates(updateData)
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
		return a.SigningKey, nil
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
	result := a.ORMInstance.First(&user, "id = ?", unmarshalled)
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
	result := a.ORMInstance.First(&user, "id = ?", claimedUser.ID)
	if result.Error != nil {
		return false, nil
	}
	user.Otp_enabled = false
	a.ORMInstance.Save(&user)
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
