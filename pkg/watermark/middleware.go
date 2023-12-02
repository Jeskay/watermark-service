package watermark

import (
	"context"
	"errors"
	"image"
	"net/http"
	authproto "watermark-service/api/v1/protos/auth"
	"watermark-service/internal"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Middleware func(Service) Service

func AuthMiddleware(authServiceAddr string) Middleware {
	return func(next Service) Service {
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
		conn, err := grpc.Dial(authServiceAddr, opts...)
		if err != nil {
			zap.L().Error("AuthMiddleware", zap.String("Dialing", "AuthService"), zap.Error(err))
			return &authMiddleware{
				next:          next,
				authAvailable: false,
			}
		}
		c := authproto.NewAuthenticationClient(conn)
		return &authMiddleware{
			next:          next,
			authClient:    c,
			authAvailable: true,
			log:           zap.L().With(zap.String("Middleware", "AuthMiddleware")),
		}
	}
}

type authMiddleware struct {
	next          Service
	authAvailable bool
	authClient    authproto.AuthenticationClient
	log           *zap.Logger
}

func (m *authMiddleware) verifyUser(ctx context.Context) (*internal.User, error) {
	token, ok := ctx.Value("token").(string)
	if !ok {
		return nil, errors.New("Empty header")
	}
	resp, err := m.authClient.VerifyJwt(ctx, &authproto.VerifyJwtRequest{Token: token})
	if err != nil {
		return nil, err
	}
	if !resp.Verified {
		return nil, errors.New("Invalid token")
	}
	var id uuid.UUID
	err = id.UnmarshalBinary(resp.User.Id)
	if err != nil {
		return nil, err
	}
	return &internal.User{
		ID:         id,
		Name:       resp.User.Name,
		Email:      resp.User.Email,
		OtpEnabled: resp.User.OtpEnabled,
	}, nil
}

func (m *authMiddleware) Add(ctx context.Context, logo image.Image, image image.Image, text string, fill bool, pos internal.Position) (string, error) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		m.log.Error("Incoming Request", zap.String("Add", "Verification"), zap.Error(err))
		return "", err
	}
	return m.next.Add(context.WithValue(ctx, "user", user), logo, image, text, fill, pos)
}

func (m *authMiddleware) Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		m.log.Error("Incoming Request", zap.String("Get", "Verification"), zap.Error(err))
		return nil, err
	}
	return m.next.Get(context.WithValue(ctx, "user", user), filters...)
}

func (m *authMiddleware) Remove(ctx context.Context, ticketID string) (int, error) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		m.log.Error("Incoming Request", zap.String("Remove", "Verification"), zap.Error(err))
		return http.StatusUnauthorized, err
	}
	return m.next.Remove(context.WithValue(ctx, "user", user), ticketID)
}

func (m *authMiddleware) ServiceStatus(ctx context.Context) (int, error) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		m.log.Error("Incoming Request", zap.String("ServiceStatus", "Verification"), zap.Error(err))
		return http.StatusUnauthorized, err
	}
	return m.next.ServiceStatus(context.WithValue(ctx, "user", user))
}
