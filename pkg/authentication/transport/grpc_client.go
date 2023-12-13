package transport

import (
	"context"
	"net/http"
	"watermark-service/api/v1/protos/auth"
	"watermark-service/internal"
	"watermark-service/internal/util"
	service "watermark-service/pkg/authentication"
	"watermark-service/pkg/authentication/endpoints"

	"github.com/go-kit/kit/endpoint"
	zapkit "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/kit/tracing/opentracing"
	grpckit "github.com/go-kit/kit/transport/grpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type grpcClient struct {
	login         endpoint.Endpoint
	register      endpoint.Endpoint
	generate      endpoint.Endpoint
	verify        endpoint.Endpoint
	verifyJwt     endpoint.Endpoint
	serviceStatus endpoint.Endpoint
}

func NewGRPCClient(conn *grpc.ClientConn) service.Service {
	logger := zapkit.NewZapSugarLogger(zap.L(), zapcore.DebugLevel)
	return &grpcClient{
		login: grpckit.NewClient(
			conn,
			auth.Authentication_ServiceDesc.ServiceName,
			"Login",
			encodeGRPCLoginRequest,
			decodeGRPCLoginResponse,
			auth.LoginResponse{},
			grpckit.ClientBefore(
				opentracing.ContextToGRPC(internal.Tracer, logger),
			),
		).Endpoint(),
		register: grpckit.NewClient(
			conn,
			auth.Authentication_ServiceDesc.ServiceName,
			"Register",
			encodeGRPCRegisterRequest,
			decodeGRPCRegisterResponse,
			auth.RegisterResponse{},
			grpckit.ClientBefore(
				opentracing.ContextToGRPC(internal.Tracer, logger),
			),
		).Endpoint(),
		generate: grpckit.NewClient(
			conn,
			auth.Authentication_ServiceDesc.ServiceName,
			"Generate",
			encodeGRPCGenerateRequest,
			decodeGRPCGenerateResponse,
			auth.GenerateResponse{},
			grpckit.ClientBefore(
				opentracing.ContextToGRPC(internal.Tracer, logger),
			),
		).Endpoint(),
		verify: grpckit.NewClient(
			conn,
			auth.Authentication_ServiceDesc.ServiceName,
			"Verify",
			encodeGRPCVerifyRequest,
			decodeGRPCVerifyResponse,
			auth.VerifyResponse{},
			grpckit.ClientBefore(
				opentracing.ContextToGRPC(internal.Tracer, logger),
			),
		).Endpoint(),
		verifyJwt: grpckit.NewClient(
			conn,
			auth.Authentication_ServiceDesc.ServiceName,
			"VerifyJwt",
			encodeGRPCVerifyJwtRequest,
			decodeGRPCVerifyJwtResponse,
			auth.VerifyJwtResponse{},
			grpckit.ClientBefore(
				opentracing.ContextToGRPC(internal.Tracer, logger),
			),
		).Endpoint(),
		serviceStatus: grpckit.NewClient(
			conn,
			auth.Authentication_ServiceDesc.ServiceName,
			"ServiceStatus",
			encodeGRPCServiceStatusRequest,
			decodeGRPCServiceStatusResponse,
			auth.ServiceStatusResponse{},
			grpckit.ClientBefore(
				opentracing.ContextToGRPC(internal.Tracer, logger),
			),
		).Endpoint(),
	}
}

func (c *grpcClient) Login(ctx context.Context, email, password string) (int64, string) {
	req := &endpoints.LoginRequest{Email: email, Password: password}
	r, err := c.login(ctx, req)
	if err != nil {
		return http.StatusUnauthorized, ""
	}
	resp := r.(*endpoints.LoginResponse)
	return resp.Status, resp.Token
}

func (c *grpcClient) Register(ctx context.Context, email, name, password string) (int32, error) {
	req := &endpoints.RegisterRequest{Email: email, Name: name, Password: password}
	r, err := c.register(ctx, req)
	if err != nil {
		return 0, err
	}
	resp := r.(*endpoints.RegisterResponse)
	return resp.UserId, util.FromString(resp.Err)
}

func (c *grpcClient) Generate(ctx context.Context) (string, string) {
	req := &endpoints.GenerateRequest{}
	r, err := c.generate(ctx, req)
	if err != nil {
		return "", ""
	}
	resp := r.(*endpoints.GenerateResponse)
	return resp.Base32, resp.OtpAuthUrl
}

func (c *grpcClient) Verify(ctx context.Context, token string) (bool, *internal.User) {
	req := &endpoints.VerifyTwoFactorRequest{Token: token}
	r, err := c.verify(ctx, req)
	if err != nil {
		return false, nil
	}
	resp := r.(*endpoints.VerifyTwoFactorResponse)
	return resp.OtpVerified, resp.User
}

func (c *grpcClient) VerifyJwt(ctx context.Context, token string) (bool, *internal.User) {
	req := &endpoints.VerifyJwtRequest{Token: token}
	r, err := c.verifyJwt(ctx, req)
	if err != nil {
		return false, nil
	}
	resp := r.(*endpoints.VerifyJwtResponse)
	return resp.Verified, resp.User
}

func (c *grpcClient) ServiceStatus(ctx context.Context) (int, error) {
	req := &endpoints.ServiceStatusRequest{}
	r, err := c.serviceStatus(ctx, req)
	if err != nil {
		return 0, err
	}
	resp := r.(*endpoints.ServiceStatusResponse)
	return resp.Code, util.FromString(resp.Err)
}

func (c *grpcClient) Validate(ctx context.Context, userId int32, token string) bool {
	return false
}

func (c *grpcClient) Disable(ctx context.Context) (bool, *internal.User) {
	return false, nil
}
