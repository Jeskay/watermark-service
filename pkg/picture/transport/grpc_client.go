package transport

import (
	"context"
	"image"
	"watermark-service/api/v1/protos/picture"
	"watermark-service/internal"
	"watermark-service/internal/util"
	service "watermark-service/pkg/picture"
	"watermark-service/pkg/picture/endpoints"

	"github.com/go-kit/kit/endpoint"
	zapkit "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/kit/tracing/opentracing"
	grpckit "github.com/go-kit/kit/transport/grpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type grpcClient struct {
	create        endpoint.Endpoint
	serviceStatus endpoint.Endpoint
}

func NewGRPCClient(conn *grpc.ClientConn) service.Service {
	logger := zapkit.NewZapSugarLogger(zap.L(), zapcore.DebugLevel)
	return &grpcClient{
		create: grpckit.NewClient(
			conn,
			"picture.Picture",
			"Create",
			encodeGRPCCreateRequest,
			decodeGRPCCreateResponse,
			picture.CreateResponse{},
			grpckit.ClientBefore(
				opentracing.ContextToGRPC(internal.Tracer, logger),
			),
		).Endpoint(),
		serviceStatus: grpckit.NewClient(
			conn,
			"picture.Picture",
			"ServiceStatus",
			encodeGRPCServiceStatusRequest,
			decodeGRPCServiceStatusResponse,
			picture.ServiceStatusResponse{},
			grpckit.ClientBefore(
				opentracing.ContextToGRPC(internal.Tracer, logger),
			),
		).Endpoint(),
	}
}

func (c *grpcClient) Create(ctx context.Context, Image image.Image, logo image.Image, text string, fill bool, pos internal.Position) (image.Image, error) {
	req := &endpoints.CreateRequest{Image: Image, Logo: logo, Text: text, Fill: fill, Pos: pos}
	r, err := c.create(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := r.(*endpoints.CreateResponse)
	return resp.Image, util.FromString(resp.Err)
}

func (c *grpcClient) ServiceStatus(ctx context.Context) (int64, error) {
	req := &endpoints.ServiceStatusRequest{}
	r, err := c.serviceStatus(ctx, req)
	if err != nil {
		return 0, err
	}
	resp := r.(*endpoints.ServiceStatusResponse)
	return resp.Code, util.FromString(resp.Err)
}
