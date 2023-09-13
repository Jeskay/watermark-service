package transport

import (
	"context"
	"watermark-service/api/v1/protos/auth"
	"watermark-service/pkg/authentication/endpoints"

	grpckit "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
)

type grpcServer struct {
	login         grpckit.Handler
	register      grpckit.Handler
	generate      grpckit.Handler
	serviceStatus grpckit.Handler
	auth.UnimplementedAuthenticationServer
}

func NewGRPCServer(ep endpoints.Set) auth.AuthenticationServer {
	return &grpcServer{
		login:         grpckit.NewServer(ep.LoginEndpoint, decodeGRPCLoginRequest, encodeGRPCLoginResponse),
		register:      grpckit.NewServer(ep.RegisterEndpoint, decodeGRPCRegisterRequest, encodeGRPCRegisterResponse),
		generate:      grpckit.NewServer(ep.GenerateEndpoint, decodeGRPCGenerateRequest, encodeGRPCGenerateResponse),
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
	req := grpcReq.(*auth.GenerateRequest)
	return endpoints.GenerateRequest{UserId: string(req.GetUserId())}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoints.ServiceStatusRequest{}, nil
}

func encodeGRPCLoginResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.LoginResponse)
	userId, err := response.User.ID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	user := auth.User{
		Id:    userId,
		Name:  response.User.Name,
		Email: response.User.Email,
	}
	return &auth.LoginResponse{Status: response.Status, User: &user}, nil
}

func encodeGRPCRegisterResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.RegisterResponse)
	user_uuid, err := uuid.Parse(response.UserId)
	if err != nil {
		return nil, err
	}
	user_id, err := user_uuid.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &auth.RegisterResponse{UserId: user_id, Error: ""}, nil
}

func encodeGRPCGenerateResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.GenerateResponse)
	return &auth.GenerateResponse{Base32: []byte(response.Base32), OtpUrl: response.OtpAuthUrl}, nil
}

func encodeGRPCServiceStatusResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.ServiceStatusResponse)
	return &auth.ServiceStatusResponse{Code: int64(response.Code), Err: response.Err}, nil
}
