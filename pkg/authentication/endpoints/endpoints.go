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
	LoginEndpoint           endpoint.Endpoint
	RegisterEndpoint        endpoint.Endpoint
	GenerateEndpoint        endpoint.Endpoint
	VerifyTwoFactorEndpoint endpoint.Endpoint
	ValidateEndpoint        endpoint.Endpoint
	VerifyJwtEndpoint       endpoint.Endpoint
	DisableEndpoint         endpoint.Endpoint
	ServiceStatusEndpoint   endpoint.Endpoint
}

func NewEndpointSet(svc authentication.Service) Set {
	return Set{
		LoginEndpoint:           MakeLoginEndpoint(svc),
		RegisterEndpoint:        MakeRegisterEndpoint(svc),
		GenerateEndpoint:        MakeGenerateEndpoint(svc),
		VerifyTwoFactorEndpoint: MakeVerifyTwoFactorEndpoint(svc),
		ValidateEndpoint:        MakeValidateEndpoint(svc),
		VerifyJwtEndpoint:       MakeVerifyJwtEndpoint(svc),
		ServiceStatusEndpoint:   MakeServiceStatusEndpoint(svc),
	}
}

func MakeLoginEndpoint(svc authentication.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		status, token := svc.Login(ctx, req.Email, req.Password)
		return LoginResponse{Status: status, Token: token}, nil
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

func MakeGenerateEndpoint(svc authentication.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GenerateRequest)
		base_32, otp_auth_url := svc.Generate(ctx, req.UserId)
		return GenerateResponse{Base32: base_32, OtpAuthUrl: otp_auth_url}, nil
	}
}

func MakeVerifyTwoFactorEndpoint(svc authentication.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyTwoFactorRequest)
		otp_verified, user := svc.Verify(ctx, req.UserId, req.Token)
		return VerifyTwoFactorResponse{OtpVerified: otp_verified, User: user}, nil
	}
}

func MakeValidateEndpoint(svc authentication.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ValidateRequest)
		otp_valid := svc.Validate(ctx, req.UserId, req.Token)
		return ValidateResponse{otp_valid}, nil
	}
}

func MakeVerifyJwtEndpoint(svc authentication.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyJwtRequest)
		verified, user := svc.VerifyJwt(ctx, req.Token)
		return VerifyJwtResponse{verified, user}, nil
	}
}

func MakeDisableEndpoint(svc authentication.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DisableRequest)
		otp_disabled, user := svc.Disable(ctx, req.UserId)
		return DisableResponse{OtpDisabled: otp_disabled, User: user}, nil
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

func (s *Set) Login(ctx context.Context, email, password string) (int64, string) {
	resp, err := s.LoginEndpoint(ctx, LoginRequest{Email: email, Password: password})
	if err != nil {
		return http.StatusUnauthorized, ""
	}
	loginResp := resp.(LoginResponse)
	return loginResp.Status, loginResp.Token
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

func (s *Set) Generate(ctx context.Context, userId string) (string, string, error) {
	resp, err := s.GenerateEndpoint(ctx, GenerateRequest{UserId: userId})
	generateResp := resp.(GenerateResponse)
	if err != nil {
		return generateResp.Base32, generateResp.OtpAuthUrl, err
	}
	return generateResp.Base32, generateResp.OtpAuthUrl, nil
}

func (s *Set) Verify(ctx context.Context, userId, token string) (bool, *internal.User) {
	resp, err := s.VerifyTwoFactorEndpoint(ctx, VerifyTwoFactorRequest{UserId: userId, Token: token})
	verifyResp := resp.(VerifyTwoFactorResponse)
	if err != nil {
		return false, verifyResp.User
	}
	return verifyResp.OtpVerified, verifyResp.User
}

func (s *Set) Validate(ctx context.Context, userId, token string) bool {
	resp, err := s.ValidateEndpoint(ctx, ValidateRequest{UserId: userId, Token: token})
	validateResp := resp.(ValidateResponse)
	if err != nil {
		return false
	}
	return validateResp.OtpValid
}

func (s *Set) VerifyJwt(ctx context.Context, token string) (bool, *internal.User) {
	resp, err := s.VerifyJwtEndpoint(ctx, VerifyJwtRequest{Token: token})
	verifyJwtResp := resp.(VerifyJwtResponse)
	if err != nil {
		return false, verifyJwtResp.User
	}
	return verifyJwtResp.Verified, verifyJwtResp.User
}

func (s *Set) Disabled(ctx context.Context, userId string) (bool, *internal.User) {
	resp, err := s.DisableEndpoint(ctx, DisableRequest{UserId: userId})
	disableResp := resp.(DisableResponse)
	if err != nil {
		return false, disableResp.User
	}
	return disableResp.OtpDisabled, disableResp.User
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
