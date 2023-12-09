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

func decodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*auth.LoginRequest)
	return endpoints.LoginRequest{Email: req.GetEmail(), Password: req.GetPassword()}, nil
}

func decodeGRPCRegisterRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*auth.RegisterRequest)
	return endpoints.RegisterRequest{Email: req.GetEmail(), Name: req.GetName(), Password: req.GetPassword()}, nil
}

func decodeGRPCGenerateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoints.GenerateRequest{}, nil
}

func decodeGRPCVerifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*auth.VerifyRequest)
	return endpoints.VerifyTwoFactorRequest{Token: req.GetToken()}, nil
}

func decodeGRPCVerifyJwtRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*auth.VerifyJwtRequest)
	return endpoints.VerifyJwtRequest{Token: req.GetToken()}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoints.ServiceStatusRequest{}, nil
}

func encodeGRPCLoginResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.LoginResponse)
	return &auth.LoginResponse{Status: response.Status, Token: response.Token}, nil
}

func encodeGRPCRegisterResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.RegisterResponse)
	return &auth.RegisterResponse{UserId: response.UserId, Error: ""}, nil
}

func encodeGRPCGenerateResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.GenerateResponse)
	return &auth.GenerateResponse{Base32: []byte(response.Base32), OtpUrl: response.OtpAuthUrl}, nil
}

func encodeGRPCVerifyResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.VerifyTwoFactorResponse)
	user := auth.User{
		Id:         response.User.ID,
		Name:       response.User.Name,
		Email:      response.User.Email,
		OtpEnabled: response.User.OtpEnabled,
	}
	return &auth.VerifyResponse{OtpVerified: response.OtpVerified, User: &user}, nil
}

func encodeGRPCVerifyJwtResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.VerifyJwtResponse)
	user := auth.User{
		Id:         response.User.ID,
		Name:       response.User.Name,
		Email:      response.User.Email,
		OtpEnabled: response.User.OtpEnabled,
	}
	return &auth.VerifyJwtResponse{Verified: response.Verified, User: &user}, nil
}

func encodeGRPCValidateResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.ValidateResponse)
	return &auth.ValidateResponse{OtpValid: response.OtpValid}, nil
}

func encodeGRPCServiceStatusResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.ServiceStatusResponse)
	return &auth.ServiceStatusResponse{Code: int64(response.Code), Err: response.Err}, nil
}
