package authentication

import (
	"context"
	"watermark-service/internal"
)

type Service interface {
	Register(ctx context.Context, email, name, password string) (string, error)
	Login(ctx context.Context, email, password string) (int64, *internal.User)
	Generate(ctx context.Context, userId string) (string, string)
	// Verify(ctx context.Context, userId []byte) (bool, internal.User)
	// Validate(ctx context.Context, userId []byte, token string) bool
	// Disable(ctx context.Context, userId []byte) (bool, internal.User)
	ServiceStatus(ctx context.Context) (int, error)
}
