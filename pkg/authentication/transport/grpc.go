package transport

import (
	"context"
	"watermark-service/api/v1/protos/auth"
	"watermark-service/internal"
	"watermark-service/pkg/authentication/endpoints"
)

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

func decodeGRPCRegisterResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*auth.RegisterResponse)
	return &endpoints.RegisterResponse{UserId: resp.GetUserId(), Err: resp.Error}, nil
}

func decodeGRPCLoginResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*auth.LoginResponse)
	return &endpoints.LoginResponse{Status: resp.GetStatus(), Token: resp.GetToken()}, nil
}

func decodeGRPCGenerateResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*auth.GenerateResponse)
	return &endpoints.GenerateResponse{Base32: string(resp.Base32), OtpAuthUrl: resp.OtpUrl}, nil
}

func decodeGRPCVerifyResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*auth.VerifyResponse)
	user := &internal.User{
		ID:         resp.User.Id,
		Name:       resp.User.Name,
		Email:      resp.User.Email,
		OtpEnabled: resp.User.OtpEnabled,
	}
	return &endpoints.VerifyTwoFactorResponse{OtpVerified: resp.OtpVerified, User: user}, nil
}

func decodeGRPCVerifyJwtResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*auth.VerifyJwtResponse)
	user := &internal.User{
		ID:         resp.User.Id,
		Name:       resp.User.Name,
		Email:      resp.User.Email,
		OtpEnabled: resp.User.OtpEnabled,
	}
	return &endpoints.VerifyJwtResponse{Verified: resp.Verified, User: user}, nil
}

func decodeGRPCServiceStatusResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*auth.ServiceStatusResponse)
	return &endpoints.ServiceStatusResponse{Code: int(resp.GetCode()), Err: resp.Err}, nil
}

func encodeGRPCRegisterRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*endpoints.RegisterRequest)
	return &auth.RegisterRequest{Email: req.Email, Name: req.Name, Password: req.Password}, nil
}

func encodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*endpoints.LoginRequest)
	return &auth.LoginRequest{Email: req.Email, Password: req.Password}, nil
}

func encodeGRPCGenerateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return &auth.GenerateRequest{}, nil
}

func encodeGRPCVerifyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*endpoints.VerifyTwoFactorRequest)
	return &auth.VerifyRequest{Token: req.Token}, nil
}

func encodeGRPCVerifyJwtRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*endpoints.VerifyJwtRequest)
	return &auth.VerifyJwtRequest{Token: req.Token}, nil
}

func encodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return &auth.ServiceStatusRequest{}, nil
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

func encodeGRPCServiceStatusResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.ServiceStatusResponse)
	return &auth.ServiceStatusResponse{Code: int64(response.Code), Err: response.Err}, nil
}
