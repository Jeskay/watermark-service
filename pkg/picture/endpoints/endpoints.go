package endpoints

import (
	"context"
	"errors"
	"image"
	"watermark-service/internal"
	"watermark-service/pkg/picture"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
)

type Set struct {
	CreateEndpoint        endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func NewEndpointSet(svc picture.Service) Set {
	return Set{
		CreateEndpoint:        MakeCreateEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
	}
}

func MakeCreateEndpoint(svc picture.Service) endpoint.Endpoint {
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		code, err := svc.Create(ctx, req.Image, req.Logo, req.Text, req.Fill, req.Pos)
		if err != nil {
			return CreateResponse{code, err.Error()}, nil
		}
		return CreateResponse{code, ""}, nil
	}
	return opentracing.TraceServer(internal.Tracer, "Create method")(endpoint)
}

func MakeServiceStatusEndpoint(svc picture.Service) endpoint.Endpoint {
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code, err := svc.ServiceStatus(ctx)
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, nil
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
	return opentracing.TraceServer(internal.Tracer, "ServiceStatus method")(endpoint)
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
