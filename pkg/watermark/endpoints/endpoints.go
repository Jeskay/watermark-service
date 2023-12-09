package endpoints

import (
	"context"
	"errors"
	"image"
	"watermark-service/internal"
	"watermark-service/pkg/watermark"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
)

type Set struct {
	GetEndpoint           endpoint.Endpoint
	AddEndpoint           endpoint.Endpoint
	RemoveEndpoint        endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func NewEndpointSet(svc watermark.Service) Set {
	return Set{
		GetEndpoint:           MakeGetEndpoint(svc),
		AddEndpoint:           MakeAddEndpoint(svc),
		RemoveEndpoint:        MakeRemoveEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
	}
}

func MakeGetEndpoint(svc watermark.Service) endpoint.Endpoint {
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		docs, err := svc.Get(ctx, req.Filters...)
		if err != nil {
			return GetResponse{Documents: docs, Err: err.Error()}, nil
		}

		return GetResponse{Documents: docs, Err: ""}, nil
	}
	return opentracing.TraceServer(internal.Tracer, "Get method")(endpoint)
}

func MakeAddEndpoint(svc watermark.Service) endpoint.Endpoint {
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		ticketID, err := svc.Add(ctx, req.Logo, req.Image, req.Text, req.Fill, req.Pos)
		if err != nil {
			return AddResponse{TicketID: ticketID, Err: err.Error()}, nil
		}
		return AddResponse{TicketID: ticketID}, nil
	}
	return opentracing.TraceServer(internal.Tracer, "Add method")(endpoint)
}

func MakeRemoveEndpoint(svc watermark.Service) endpoint.Endpoint {
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RemoveRequest)
		code, err := svc.Remove(ctx, req.TicketID)
		if err != nil {
			return RemoveResponse{Code: code, Err: err.Error()}, nil
		}
		return RemoveResponse{Code: code, Err: ""}, nil
	}
	return opentracing.TraceServer(internal.Tracer, "Remove method")(endpoint)
}

func MakeServiceStatusEndpoint(svc watermark.Service) endpoint.Endpoint {
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

func (s *Set) Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	resp, err := s.GetEndpoint(ctx, GetRequest{Filters: filters})
	if err != nil {
		return []internal.Document{}, err
	}
	getResp := resp.(GetResponse)
	if getResp.Err != "" {
		return []internal.Document{}, errors.New(getResp.Err)
	}
	return getResp.Documents, nil
}

func (s *Set) Add(ctx context.Context, logo image.Image, image image.Image, text string, fill bool, pos internal.Position) (string, error) {
	resp, err := s.AddEndpoint(ctx, AddRequest{Logo: logo, Image: image, Text: text, Fill: fill, Pos: pos})
	if err != nil {
		return "", err
	}
	addResp := resp.(AddResponse)
	if addResp.Err != "" {
		return "", errors.New(addResp.Err)
	}
	return addResp.TicketID, nil
}

func (s *Set) Remove(ctx context.Context, ticketID string) (int, error) {
	resp, err := s.RemoveEndpoint(ctx, RemoveRequest{TicketID: ticketID})
	removeResp := resp.(RemoveResponse)
	if err != nil {
		return removeResp.Code, err
	}
	if removeResp.Err != "" {
		return removeResp.Code, errors.New(removeResp.Err)
	}
	return removeResp.Code, nil
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
