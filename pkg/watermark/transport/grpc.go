package transport

import (
	"context"
	"watermark-service/api/v1/protos/watermark"
	"watermark-service/pkg/watermark/endpoints"

	grpckit "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	create        grpckit.Handler
	serviceStatus grpckit.Handler
	watermark.UnimplementedWatermarkServer
}

func NewGRPCServer(ep endpoints.Set) watermark.WatermarkServer {
	return &grpcServer{
		create: grpckit.NewServer(
			ep.CreateEndpoint,
			decodeGRPCCreateRequest,
			encodeGRPCCreateResponse,
		),
		serviceStatus: grpckit.NewServer(
			ep.ServiceStatusEndpoint,
			decodeGRPCServiceStatusRequest,
			encodeGRPCServiceStatusResponse,
		),
	}
}

func (g *grpcServer) Create(ctx context.Context, r *watermark.CreateRequest) (*watermark.CreateResponse, error) {
	_, rep, err := g.create.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.CreateResponse), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *watermark.ServiceStatusRequest) (*watermark.ServiceStatusResponse, error) {
	_, rep, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.ServiceStatusResponse), nil
}

func decodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.CreateRequest)
	return endpoints.CreateRequest{Logo: nil, Text: req.Text, Fill: req.Fill}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoints.ServiceStatusRequest{}, nil
}

func encodeGRPCCreateResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(*watermark.CreateResponse)
	return endpoints.CreateResponse{Image: nil, Err: response.GetErr()}, nil
}

func encodeGRPCServiceStatusResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(*watermark.ServiceStatusResponse)
	return endpoints.ServiceStatusResponse{Code: response.GetCode(), Err: response.GetErr()}, nil
}
