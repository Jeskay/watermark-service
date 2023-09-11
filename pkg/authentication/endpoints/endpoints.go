package endpoints

import (
	"context"
	"errors"
	"net/http"
	"watermark-service/internal"
	"watermark-service/pkg/authentication"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	LoginEndpoint         endpoint.Endpoint
	RegisterEndpoint      endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func NewEndpointSet(svc authentication.Service) Set {
	return Set{
		LoginEndpoint:         MakeLoginEndpoint(svc),
		RegisterEndpoint:      MakeRegisterEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
	}
}

func MakeLoginEndpoint(svc authentication.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		status, user := svc.Login(ctx, req.Email, req.Password)
		return LoginResponse{Status: status, User: user}, nil
	}
}

func MakeRegisterEndpoint(svc authentication.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RegisterRequest)
		userId, err := svc.Register(ctx, req.Email, req.Name, req.Password)
		if err != nil {
			return RegisterResponse{UserId: userId, Err: err.Error()}, nil
		}
		return RegisterResponse{UserId: userId, Err: ""}, nil
	}
}

func MakeServiceStatusEndpoint(svc authentication.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code, err := svc.ServiceStatus(ctx)
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, nil
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
}

func (s *Set) Login(ctx context.Context, email, password string) (int64, *internal.User) {
	resp, err := s.LoginEndpoint(ctx, LoginRequest{Email: email, Password: password})
	if err != nil {
		return http.StatusUnauthorized, nil
	}
	loginResp := resp.(LoginResponse)
	return loginResp.Status, loginResp.User
}

func (s *Set) Register(ctx context.Context, email, name, password string) (string, error) {
	resp, err := s.RegisterEndpoint(ctx, RegisterRequest{Email: email, Name: name, Password: password})
	registerResp := resp.(RegisterResponse)
	if err != nil {
		return registerResp.UserId, err
	}
	if registerResp.Err != "" {
		return registerResp.UserId, errors.New(registerResp.Err)
	}
	return registerResp.UserId, nil
}

func (s *Set) ServiceStatus(ctx context.Context) (int, error) {
	resp, err := s.ServiceStatusEndpoint(ctx, ServiceStatusRequest{})
	svcStatusResp := resp.(ServiceStatusResponse)
	if err != nil {
		return svcStatusResp.Code, err
	}
	if svcStatusResp.Err != "" {
		return svcStatusResp.Code, errors.New(svcStatusResp.Err)
	}
	return svcStatusResp.Code, nil
}
