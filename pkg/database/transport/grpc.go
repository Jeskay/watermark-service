package transport

import (
	"context"
	"watermark-service/api/v1/protos/db"
	"watermark-service/internal"
	"watermark-service/pkg/database/endpoints"

	grpckit "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	get           grpckit.Handler
	update        grpckit.Handler
	add           grpckit.Handler
	remove        grpckit.Handler
	serviceStatus grpckit.Handler
	db.UnimplementedDatabaseServer
}

func NewGRPCServer(ep endpoints.Set) db.DatabaseServer {
	return &grpcServer{
		get:           grpckit.NewServer(ep.GetEndpoint, decodeGRPCGetRequest, encodeGRPCGetResponse),
		update:        grpckit.NewServer(ep.UpdateEndpoint, decodeGRPCUpdateRequest, encodeGRPCUpdateResponse),
		add:           grpckit.NewServer(ep.AddEndpoint, decodeGRPCAddRequest, encodeGRPCAddResponse),
		remove:        grpckit.NewServer(ep.RemoveEndpoint, decodeGRPCRemoveRequest, encodeGRPCRemoveResponse),
		serviceStatus: grpckit.NewServer(ep.ServiceStatusEndpoint, decodeGRPCServiceStatusRequest, encodeGRPCServiceStatusResponse),
	}
}

func (g *grpcServer) Get(ctx context.Context, r *db.GetRequest) (*db.GetResponse, error) {
	_, resp, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*db.GetResponse), nil
}

func (g *grpcServer) Update(ctx context.Context, r *db.UpdateRequest) (*db.UpdateResponse, error) {
	_, resp, err := g.update.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*db.UpdateResponse), nil
}

func (g *grpcServer) Add(ctx context.Context, r *db.AddRequest) (*db.AddResponse, error) {
	_, resp, err := g.update.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*db.AddResponse), nil
}

func (g *grpcServer) Remove(ctx context.Context, r *db.RemoveRequest) (*db.RemoveResponse, error) {
	_, resp, err := g.update.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*db.RemoveResponse), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *db.ServiceStatusRequest) (*db.ServiceStatusResponse, error) {
	_, resp, err := g.update.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*db.ServiceStatusResponse), nil
}

func decodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*db.GetRequest)
	var filters []internal.Filter
	for _, f := range req.Filters {
		filters = append(filters, internal.Filter{Key: f.Key, Value: f.Value})
	}
	return endpoints.GetRequest{Filters: filters}, nil
}

func decodeGRPCUpdateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*db.UpdateRequest)
	doc := &internal.Document{
		Content:   req.Document.GetContent(),
		Title:     req.Document.GetTitle(),
		Author:    req.Document.GetAuthor(),
		Topic:     req.Document.GetTopic(),
		Watermark: req.Document.GetWatermark(),
	}
	return endpoints.UpdateRequest{TicketID: req.TicketID, Document: doc}, nil
}

func decodeGRPCRemoveRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*db.RemoveRequest)
	return endpoints.RemoveRequest{TicketID: req.TicketID}, nil
}

func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*db.AddRequest)
	doc := &internal.Document{
		Content:   req.Document.GetContent(),
		Title:     req.Document.GetTitle(),
		Author:    req.Document.GetAuthor(),
		Topic:     req.Document.GetTopic(),
		Watermark: req.Document.GetWatermark(),
	}
	return endpoints.AddRequest{Document: doc}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoints.ServiceStatusRequest{}, nil
}

func encodeGRPCGetResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*db.GetResponse)
	var docs []internal.Document
	for _, d := range response.Documents {
		doc := internal.Document{
			Content:   d.GetContent(),
			Title:     d.GetTitle(),
			Author:    d.GetAuthor(),
			Topic:     d.GetAuthor(),
			Watermark: d.GetWatermark(),
		}
		docs = append(docs, doc)
	}
	return endpoints.GetResponse{Documents: docs, Err: response.GetErr()}, nil
}

func encodeGRPCUpdateResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*db.UpdateResponse)
	return endpoints.UpdateResponse{Code: int(response.GetCode()), Err: response.GetErr()}, nil
}

func encodeGRPCRemoveResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*db.RemoveResponse)
	return endpoints.RemoveResponse{Code: int(response.GetCode()), Err: response.GetErr()}, nil
}

func encodeGRPCServiceStatusResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*db.ServiceStatusResponse)
	return endpoints.ServiceStatusResponse{Code: int(response.GetCode()), Err: response.GetErr()}, nil
}

func encodeGRPCAddResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(*db.AddResponse)
	return endpoints.AddResponse{TicketID: response.GetTicketID(), Err: response.GetErr()}, nil
}
