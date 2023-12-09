package transport

import (
	"context"
	"watermark-service/api/v1/protos/watermark"
	"watermark-service/internal"
	"watermark-service/internal/util"
	"watermark-service/pkg/watermark/endpoints"

	zapkit "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/kit/tracing/opentracing"
	grpckit "github.com/go-kit/kit/transport/grpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type grpcServer struct {
	get           grpckit.Handler
	add           grpckit.Handler
	remove        grpckit.Handler
	serviceStatus grpckit.Handler
	watermark.UnimplementedWatermarkServer
}

func NewGRPCServer(ep endpoints.Set) watermark.WatermarkServer {
	logger := zapkit.NewZapSugarLogger(zap.L(), zapcore.DebugLevel)
	return &grpcServer{
		get: grpckit.NewServer(
			ep.GetEndpoint,
			decodeGRPCGetRequest,
			encodeGRPCGetResponse,
			grpckit.ServerBefore(
				opentracing.GRPCToContext(internal.Tracer, "Get method", logger),
			),
		),
		add: grpckit.NewServer(
			ep.AddEndpoint,
			decodeGRPCAddRequest,
			encodeGRPCAddResponse,
			grpckit.ServerBefore(
				opentracing.GRPCToContext(internal.Tracer, "Add method", zapkit.NewZapSugarLogger(zap.L(), zapcore.DebugLevel)),
			),
		),
		remove: grpckit.NewServer(
			ep.RemoveEndpoint,
			decodeGRPCRemoveRequest,
			encodeGRPCRemoveResponse,
			grpckit.ServerBefore(
				opentracing.GRPCToContext(internal.Tracer, "Remove method", zapkit.NewZapSugarLogger(zap.L(), zapcore.DebugLevel)),
			),
		),
		serviceStatus: grpckit.NewServer(
			ep.ServiceStatusEndpoint,
			decodeGRPCServiceStatusRequest,
			encodeGRPCServiceStatusResponse,
			grpckit.ServerBefore(
				opentracing.GRPCToContext(internal.Tracer, "ServiceStatus method", zapkit.NewZapSugarLogger(zap.L(), zapcore.DebugLevel)),
			),
		),
	}
}

func (g *grpcServer) Get(ctx context.Context, r *watermark.GetRequest) (*watermark.GetResponse, error) {
	_, resp, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*watermark.GetResponse), nil
}

func (g *grpcServer) Add(ctx context.Context, r *watermark.AddRequest) (*watermark.AddResponse, error) {
	_, resp, err := g.add.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*watermark.AddResponse), nil
}

func (g *grpcServer) Remove(ctx context.Context, r *watermark.RemoveRequest) (*watermark.RemoveResponse, error) {
	_, resp, err := g.remove.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*watermark.RemoveResponse), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *watermark.ServiceStatusRequest) (*watermark.ServiceStatusResponse, error) {
	_, resp, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*watermark.ServiceStatusResponse), nil
}

func decodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.GetRequest)
	var filters []internal.Filter
	for _, f := range req.Filters {
		filters = append(filters, internal.Filter{Key: f.Key, Value: f.Value})
	}
	return endpoints.GetRequest{Filters: filters}, nil
}

func decodeGRPCRemoveRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.RemoveRequest)
	return endpoints.RemoveRequest{TicketID: req.TicketID}, nil
}

func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.AddRequest)
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
	var docs []*watermark.Document
	for _, d := range response.Documents {
		ticket_id, err := d.ID.MarshalBinary()
		if err != nil {
			return nil, err
		}
		author_id, err := d.AuthorId.MarshalBinary()
		if err != nil {
			return nil, err
		}
		doc := watermark.Document{
			TicketId: ticket_id,
			AuthorId: author_id,
			Title:    d.Title,
			ImageUrl: d.ImageUrl,
		}
		docs = append(docs, &doc)
	}
	return &watermark.GetResponse{Documents: docs, Err: response.Err}, nil
}

func encodeGRPCRemoveResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.RemoveResponse)
	return &watermark.RemoveResponse{Code: int64(response.Code), Err: response.Err}, nil
}

func encodeGRPCServiceStatusResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.ServiceStatusResponse)
	return &watermark.ServiceStatusResponse{Code: int64(response.Code), Err: response.Err}, nil
}

func encodeGRPCAddResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	response := grpcResponse.(endpoints.AddResponse)
	return &watermark.AddResponse{TicketID: response.TicketID, Err: response.Err}, nil
}
