package transport

import (
	"context"
	"watermark-service/api/v1/protos/picture"
	"watermark-service/internal"
	"watermark-service/pkg/picture/endpoints"

	zapkit "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/kit/tracing/opentracing"
	grpckit "github.com/go-kit/kit/transport/grpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type grpcServer struct {
	create        grpckit.Handler
	serviceStatus grpckit.Handler
	picture.UnimplementedPictureServer
}

func NewGRPCServer(ep endpoints.Set) picture.PictureServer {
	logger := zapkit.NewZapSugarLogger(zap.L(), zapcore.DebugLevel)
	return &grpcServer{
		create: grpckit.NewServer(
			ep.CreateEndpoint,
			decodeGRPCCreateRequest,
			encodeGRPCCreateResponse,
			grpckit.ServerBefore(
				opentracing.GRPCToContext(
					internal.Tracer,
					"Create method",
					logger,
				),
			),
		),
		serviceStatus: grpckit.NewServer(
			ep.ServiceStatusEndpoint,
			decodeGRPCServiceStatusRequest,
			encodeGRPCServiceStatusResponse,
			grpckit.ServerBefore(
				opentracing.GRPCToContext(
					internal.Tracer,
					"ServiceStatus method",
					logger,
				),
			),
		),
	}
}

func (g *grpcServer) Create(ctx context.Context, r *picture.CreateRequest) (*picture.CreateResponse, error) {
	_, rep, err := g.create.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*picture.CreateResponse), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *picture.ServiceStatusRequest) (*picture.ServiceStatusResponse, error) {
	_, rep, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*picture.ServiceStatusResponse), nil
}
