package authentication

import (
	"context"
	"watermark-service/internal"
)

type Service interface {
	Register(ctx context.Context, email, name, password string) (string, error)
	Login(ctx context.Context, email, password string) (int64, string)
	Generate(ctx context.Context) (string, string)
	Verify(ctx context.Context, token string) (bool, *internal.User)
	Validate(ctx context.Context, userId, token string) bool
	VerifyJwt(ctx context.Context, token string) (bool, *internal.User)
	Disable(ctx context.Context) (bool, *internal.User)
	ServiceStatus(ctx context.Context) (int, error)
}
