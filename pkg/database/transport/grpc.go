package transport

import (
	"context"
	"watermark-service/api/v1/protos/db"
	"watermark-service/internal"
	"watermark-service/internal/util"
	"watermark-service/pkg/database/endpoints"

	grpckit "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	get           grpckit.Handler
	add           grpckit.Handler
	remove        grpckit.Handler
	serviceStatus grpckit.Handler
	db.UnimplementedDatabaseServer
}

func NewGRPCServer(ep endpoints.Set) db.DatabaseServer {
	return &grpcServer{
		get:           grpckit.NewServer(ep.GetEndpoint, decodeGRPCGetRequest, encodeGRPCGetResponse),
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

func (g *grpcServer) Add(ctx context.Context, r *db.AddRequest) (*db.AddResponse, error) {
	_, resp, err := g.add.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*db.AddResponse), nil
}

func (g *grpcServer) Remove(ctx context.Context, r *db.RemoveRequest) (*db.RemoveResponse, error) {
	_, resp, err := g.remove.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*db.RemoveResponse), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *db.ServiceStatusRequest) (*db.ServiceStatusResponse, error) {
	_, resp, err := g.serviceStatus.ServeGRPC(ctx, r)
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

func decodeGRPCRemoveRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*db.RemoveRequest)
	return endpoints.RemoveRequest{TicketID: req.TicketID}, nil
}

func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*db.AddRequest)
	return endpoints.AddRequest{
		Logo:  util.ByteToImage(req.Logo.Data, req.Logo.Type),
		Image: util.ByteToImage(req.Image.Data, req.Image.Type),
		Text:  req.Text,
		Fill:  req.Fill,
		Pos:   internal.PositionFromString(req.Pos.String()),
	}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoints.ServiceStatusRequest{}, nil
}

func encodeGRPCGetResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.GetResponse)
	var docs []*db.Document
	for _, d := range response.Documents {
		ticket_id, err := d.ID.MarshalBinary()
		if err != nil {
			return nil, err
		}
		author_id, err := d.AuthorId.MarshalBinary()
		if err != nil {
			return nil, err
		}
		doc := db.Document{
			TicketId: ticket_id,
			AuthorId: author_id,
			Title:    d.Title,
			ImageUrl: d.ImageUrl,
		}
		docs = append(docs, &doc)
	}
	return &db.GetResponse{Documents: docs, Err: response.Err}, nil
}

func encodeGRPCRemoveResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.RemoveResponse)
	return &db.RemoveResponse{Code: int64(response.Code), Err: response.Err}, nil
}

func encodeGRPCServiceStatusResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.ServiceStatusResponse)
	return &db.ServiceStatusResponse{Code: int64(response.Code), Err: response.Err}, nil
}

func encodeGRPCAddResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.AddResponse)
	return &db.AddResponse{TicketID: response.TicketID, Err: response.Err}, nil
}
