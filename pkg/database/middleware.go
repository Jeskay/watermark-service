package database

import (
	"context"
	"errors"
	"log"
	"net/http"
	authproto "watermark-service/api/v1/protos/auth"
	"watermark-service/internal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Middleware func(Service) Service

func AuthMiddleware(authServiceAddr string) Middleware {
	return func(next Service) Service {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn, err := grpc.Dial(authServiceAddr, opts...)
		if err != nil {
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
		}
	}
}

type authMiddleware struct {
	next          Service
	authAvailable bool
	authClient    authproto.AuthenticationClient
}

func (m *authMiddleware) verifyUser(ctx context.Context) (*authproto.User, error) {
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
	return resp.User, nil
}

func (m *authMiddleware) Add(ctx context.Context, doc *internal.Document) (int64, error) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		log.Println(err)
		return http.StatusUnauthorized, err
	}
	return m.next.Add(context.WithValue(ctx, "user", user), doc)
}

func (m *authMiddleware) Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return m.next.Get(context.WithValue(ctx, "user", user), filters...)
}

func (m *authMiddleware) Update(ctx context.Context, ticketID int64, doc *internal.Document) (int, error) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		log.Println(err)
		return http.StatusUnauthorized, err
	}
	return m.next.Update(context.WithValue(ctx, "user", user), ticketID, doc)
}

func (m *authMiddleware) Remove(ctx context.Context, ticketID int64) (int, error) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		log.Println(err)
		return http.StatusUnauthorized, err
	}
	return m.next.Remove(context.WithValue(ctx, "user", user), ticketID)
}

func (m *authMiddleware) ServiceStatus(ctx context.Context) (int, error) {
	user, err := m.verifyUser(ctx)
	if err != nil {
		log.Println(err)
		return http.StatusUnauthorized, err
	}
	return m.next.ServiceStatus(context.WithValue(ctx, "user", user))
}
