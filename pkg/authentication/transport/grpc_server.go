package transport

import (
	"context"
	"watermark-service/api/v1/protos/auth"
	"watermark-service/pkg/authentication/endpoints"

	grpckit "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	login         grpckit.Handler
	register      grpckit.Handler
	generate      grpckit.Handler
	verify        grpckit.Handler
	verifyJwt     grpckit.Handler
	validate      grpckit.Handler
	serviceStatus grpckit.Handler
	auth.UnimplementedAuthenticationServer
}

func NewGRPCServer(ep endpoints.Set) auth.AuthenticationServer {
	return &grpcServer{
		login:         grpckit.NewServer(ep.LoginEndpoint, decodeGRPCLoginRequest, encodeGRPCLoginResponse),
		register:      grpckit.NewServer(ep.RegisterEndpoint, decodeGRPCRegisterRequest, encodeGRPCRegisterResponse),
		generate:      grpckit.NewServer(ep.GenerateEndpoint, decodeGRPCGenerateRequest, encodeGRPCGenerateResponse),
		verify:        grpckit.NewServer(ep.VerifyTwoFactorEndpoint, decodeGRPCVerifyRequest, encodeGRPCVerifyResponse),
		verifyJwt:     grpckit.NewServer(ep.VerifyJwtEndpoint, decodeGRPCVerifyJwtRequest, encodeGRPCVerifyJwtResponse),
		serviceStatus: grpckit.NewServer(ep.ServiceStatusEndpoint, decodeGRPCServiceStatusRequest, encodeGRPCServiceStatusResponse),
	}
}

func (g *grpcServer) Login(ctx context.Context, r *auth.LoginRequest) (*auth.LoginResponse, error) {
	_, resp, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*auth.LoginResponse), nil
}

func (g *grpcServer) Register(ctx context.Context, r *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	_, resp, err := g.register.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*auth.RegisterResponse), nil
}

func (g *grpcServer) Generate(ctx context.Context, r *auth.GenerateRequest) (*auth.GenerateResponse, error) {
	_, resp, err := g.generate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*auth.GenerateResponse), nil
}

func (g *grpcServer) Verify(ctx context.Context, r *auth.VerifyRequest) (*auth.VerifyResponse, error) {
	_, resp, err := g.verify.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*auth.VerifyResponse), nil
}

func (g *grpcServer) VerifyJwt(ctx context.Context, r *auth.VerifyJwtRequest) (*auth.VerifyJwtResponse, error) {
	_, resp, err := g.verifyJwt.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*auth.VerifyJwtResponse), nil
}

func (g *grpcServer) Validate(ctx context.Context, r *auth.ValidateRequest) (*auth.ValidateResponse, error) {
	_, resp, err := g.validate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*auth.ValidateResponse), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *auth.ServiceStatusRequest) (*auth.ServiceStatusResponse, error) {
	_, resp, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*auth.ServiceStatusResponse), nil
}
