package authentication

import (
	"context"
	"errors"
	"watermark-service/internal"

	"go.uber.org/zap"
)

type Middleware func(Service) Service

func AuthMiddleware() Middleware {
	return func(next Service) Service {
		return &authMiddleware{
			next: next,
			log:  zap.L().With(zap.String("Middleware", "AuthenticationMiddleware")),
		}

	}
}

type authMiddleware struct {
	next Service
	log  *zap.Logger
}

func (m *authMiddleware) verifyUser(ctx context.Context) (*internal.User, error) {
	token, ok := ctx.Value("token").(string)
	if !ok {
		return nil, errors.New("Empty header")
	}
	ok, user := m.next.VerifyJwt(ctx, token)
	if !ok {
		return nil, errors.New("Invalid token")
	}
	return user, nil
}

func (m *authMiddleware) Generate(ctx context.Context) (string, string) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		m.log.Error("Generate", zap.Error(err))
		return "", ""
	}
	return m.next.Generate(context.WithValue(ctx, "user", user))
}

func (m *authMiddleware) Verify(ctx context.Context, token string) (bool, *internal.User) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		m.log.Error("Generate", zap.Error(err))
		return false, nil
	}
	return m.next.Verify(context.WithValue(ctx, "user", user), token)
}

func (m *authMiddleware) Disable(ctx context.Context) (bool, *internal.User) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		m.log.Error("Generate", zap.Error(err))
		return false, nil
	}
	return m.next.Disable(context.WithValue(ctx, "user", user))
}

func (m *authMiddleware) Validate(ctx context.Context, userId int32, token string) bool {
	return m.next.Validate(ctx, userId, token)
}

func (m *authMiddleware) Login(ctx context.Context, email, password string) (int64, string) {
	if token2FA := ctx.Value("2FA"); token2FA != nil {
		return m.next.Login(context.WithValue(ctx, "2FA", token2FA.(string)), email, password)
	}
	return m.next.Login(ctx, email, password)
}

func (m *authMiddleware) Register(ctx context.Context, email, name, password string) (int32, error) {
	return m.next.Register(ctx, email, name, password)
}

func (m *authMiddleware) ServiceStatus(ctx context.Context) (int, error) {
	return m.next.ServiceStatus(ctx)
}

func (m *authMiddleware) VerifyJwt(ctx context.Context, token string) (bool, *internal.User) {
	return m.next.VerifyJwt(ctx, token)
}
