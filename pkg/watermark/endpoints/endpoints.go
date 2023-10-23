package endpoints

import (
	"context"
	"errors"
	"image"
	"os"
	"watermark-service/pkg/watermark"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	CreateEndpoint        endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func NewEndpointSet(svc watermark.Service) Set {
	return Set{
		CreateEndpoint:        MakeCreateEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
	}
}

func MakeCreateEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		code, err := svc.Create(ctx, req.Image, req.Logo, req.Text, req.Fill, req.Pos)
		if err != nil {
			return CreateResponse{code, err.Error()}, nil
		}
		return CreateResponse{code, ""}, nil
	}
}

func MakeServiceStatusEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code, err := svc.ServiceStatus(ctx)
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, nil
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
}

func (s *Set) Create(ctx context.Context) (image.Image, error) {
	resp, err := s.CreateEndpoint(ctx, CreateRequest{})
	if err != nil {
		return nil, err
	}
	createResp := resp.(CreateResponse)
	if createResp.Err != "" {
		return nil, errors.New(createResp.Err)
	}
	return createResp.Image, nil
}

func (s *Set) ServiceStatus(ctx context.Context) (int64, error) {
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

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
