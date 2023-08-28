package endpoints

import (
	"context"
	"errors"
	"os"
	"watermark-service/internal"
	"watermark-service/pkg/watermark"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	GetEndpoint           endpoint.Endpoint
	AddDocumentEndpoint   endpoint.Endpoint
	StatusEndpoint        endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
	WatermarkEndpoint     endpoint.Endpoint
}

func NewEndpointSet(svc watermark.Service) Set {
	return Set{
		GetEndpoint:           MakeGetEndpoint(svc),
		AddDocumentEndpoint:   MakeAddDocumentEndpoint(svc),
		StatusEndpoint:        MakeStatusEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
		WatermarkEndpoint:     MakeWatermarkEndpoint(svc),
	}
}

func MakeGetEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		docs, err := svc.Get(ctx, req.Filters...)
		if err != nil {
			return GetResponse{docs, err.Error()}, nil
		}
		return GetResponse{docs, ""}, nil
	}
}

func MakeAddDocumentEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddDocumentRequest)
		ticketID, err := svc.AddDocument(ctx, req.Document)
		if err != nil {
			return AddDocumentResponse{TicketID: ticketID, Err: err.Error()}, nil
		}
		return AddDocumentResponse{TicketID: ticketID, Err: ""}, nil
	}
}

func MakeStatusEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StatusRequest)
		status, err := svc.Status(ctx, req.TicketID)
		if err != nil {
			return StatusResponse{Status: status, Err: err.Error()}, nil
		}
		return StatusResponse{Status: status, Err: ""}, nil
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

func MakeWatermarkEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(WatermarkRequest)
		code, err := svc.Watermark(ctx, req.TicketID, req.Mark)
		if err != nil {
			return WatermarkResponse{Code: code, Err: err.Error()}, nil
		}
		return WatermarkResponse{Code: code, Err: ""}, nil
	}
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

func (s *Set) ServiceStatus(ctx context.Context) (int, error) {
	resp, err := s.StatusEndpoint(ctx, StatusRequest{})
	svcStatusResp := resp.(ServiceStatusResponse)
	if err != nil {
		return svcStatusResp.Code, err
	}
	if svcStatusResp.Err != "" {
		return svcStatusResp.Code, errors.New(svcStatusResp.Err)
	}
	return svcStatusResp.Code, nil
}

func (s *Set) AddDocument(ctx context.Context, doc *internal.Document) (int64, error) {
	resp, err := s.AddDocumentEndpoint(ctx, AddDocumentRequest{Document: doc})
	if err != nil {
		return -1, err
	}
	addResp := resp.(AddDocumentResponse)
	if addResp.Err != "" {
		return -1, errors.New(addResp.Err)
	}
	return addResp.TicketID, nil
}

func (s *Set) Status(ctx context.Context, ticketID int64) (internal.Status, error) {
	resp, err := s.StatusEndpoint(ctx, StatusRequest{TicketID: ticketID})
	if err != nil {
		return internal.Failed, err
	}
	statusResp := resp.(StatusResponse)
	if statusResp.Err != "" {
		return internal.Failed, errors.New(statusResp.Err)
	}
	return statusResp.Status, nil
}

func (s *Set) Watermark(ctx context.Context, ticketID int64, mark string) (int, error) {
	resp, err := s.WatermarkEndpoint(ctx, WatermarkRequest{TicketID: ticketID, Mark: mark})
	watermarkResp := resp.(WatermarkResponse)
	if err != nil {
		return watermarkResp.Code, err
	}
	if watermarkResp.Err != "" {
		return watermarkResp.Code, errors.New(watermarkResp.Err)
	}
	return watermarkResp.Code, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
