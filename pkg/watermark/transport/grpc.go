package transport

import (
	"context"
	"watermark-service/api/v1/protos/watermark"
	"watermark-service/internal"
	"watermark-service/pkg/watermark/endpoints"

	grpckit "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	get           grpckit.Handler
	status        grpckit.Handler
	addDocument   grpckit.Handler
	watermark     grpckit.Handler
	serviceStatus grpckit.Handler
	watermark.UnimplementedWatermarkServer
}

func NewGRPCServer(ep endpoints.Set) watermark.WatermarkServer {
	return &grpcServer{
		get: grpckit.NewServer(
			ep.GetEndpoint,
			decodeGRPCGetRequest,
			encodeGRPCGetResponse,
		),
		status: grpckit.NewServer(
			ep.StatusEndpoint,
			decodeGRPCStatusRequest,
			encodeGRPCStatusResponse,
		),
		addDocument: grpckit.NewServer(
			ep.AddDocumentEndpoint,
			decodeGRPCAddDocumentRequest,
			encodeGRPCAddDocumentResponse,
		),
		watermark: grpckit.NewServer(
			ep.WatermarkEndpoint,
			decodeGRPCWatermarkRequest,
			encodeGRPCWatermarkResponse,
		),
		serviceStatus: grpckit.NewServer(
			ep.ServiceStatusEndpoint,
			decodeGRPCServiceStatusRequest,
			encodeGRPCStatusResponse,
		),
	}
}

func (g *grpcServer) Get(ctx context.Context, r *watermark.GetRequest) (*watermark.GetResponse, error) {
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.GetResponse), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *watermark.ServiceStatusRequest) (*watermark.ServiceStatusResponse, error) {
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.ServiceStatusResponse), nil
}

func (g *grpcServer) AddDocument(ctx context.Context, r *watermark.AddDocumentRequest) (*watermark.AddDocumentResponse, error) {
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.AddDocumentResponse), nil
}

func (g *grpcServer) Status(ctx context.Context, r *watermark.StatusRequest) (*watermark.StatusResponse, error) {
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.StatusResponse), nil
}

func (g *grpcServer) Watermark(ctx context.Context, r *watermark.WatermarkRequest) (*watermark.WatermarkResponse, error) {
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.WatermarkResponse), nil
}

func decodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.GetRequest)
	var filters []internal.Filter
	for _, f := range req.Filters {
		filters = append(filters, internal.Filter{Key: f.GetKey(), Value: f.GetValue()})
	}
	return endpoints.GetRequest{Filters: filters}, nil
}

func decodeGRPCStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.StatusRequest)
	return endpoints.StatusRequest{TicketID: req.GetTicketID()}, nil
}

func decodeGRPCWatermarkRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.WatermarkRequest)
	return endpoints.WatermarkRequest{TicketID: req.GetTicketID(), Mark: req.GetMark()}, nil
}

func decodeGRPCAddDocumentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.AddDocumentRequest)
	doc := &internal.Document{
		Content:   req.Document.GetContent(),
		Title:     req.Document.GetTitle(),
		Author:    req.Document.GetAuthor(),
		Topic:     req.Document.GetTopic(),
		Watermark: req.Document.GetWatermark(),
	}
	return endpoints.AddDocumentRequest{Document: doc}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoints.ServiceStatusRequest{}, nil
}

func encodeGRPCGetResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(*watermark.GetResponse)
	var docs []internal.Document
	for _, d := range response.Documents {
		doc := internal.Document{
			Content:   d.GetContent(),
			Title:     d.GetTitle(),
			Author:    d.GetAuthor(),
			Topic:     d.GetTopic(),
			Watermark: d.GetWatermark(),
		}
		docs = append(docs, doc)
	}
	return endpoints.GetResponse{Documents: docs}, nil
}

func encodeGRPCStatusResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(*watermark.StatusResponse)
	return endpoints.StatusResponse{Status: internal.Status(response.GetStatus()), Err: response.GetErr()}, nil
}

func encodeGRPCWatermarkResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(*watermark.WatermarkResponse)
	return endpoints.WatermarkResponse{Code: int(response.GetCode()), Err: response.GetErr()}, nil
}

func encodeGRPCAddDocumentResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(*watermark.AddDocumentResponse)
	return endpoints.AddDocumentResponse{TicketID: response.GetTicketID(), Err: response.GetErr()}, nil
}

func encodeGRPCServiceStatusResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(*watermark.ServiceStatusResponse)
	return endpoints.ServiceStatusResponse{Code: int(response.GetCode()), Err: response.GetErr()}, nil
}
