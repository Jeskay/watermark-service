package transport

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
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

type grpcServer struct {
	create        grpckit.Handler
	serviceStatus grpckit.Handler
	picture.UnimplementedPictureServer
}

type grpcClient struct {
	create        endpoint.Endpoint
	serviceStatus endpoint.Endpoint
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

func decodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	var Logo image.Image
	req := grpcReq.(*picture.CreateRequest)
	img := req.GetImage()
	logo := req.GetLogo()
	if logo != nil {
		Logo = getImageFromByte(logo.Data, logo.Type)
	}
	pos := internal.PositionFromString(req.Pos.String())
	image := getImageFromByte(img.Data, img.Type)
	return endpoints.CreateRequest{Image: image, Logo: Logo, Text: req.Text, Fill: req.Fill, Pos: pos}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoints.ServiceStatusRequest{}, nil
}

func encodeGRPCCreateResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(endpoints.CreateResponse)
	buf := new(bytes.Buffer)
	png.Encode(buf, response.Image)
	return &picture.CreateResponse{Image: buf.Bytes(), Err: response.Err}, nil
}

func encodeGRPCServiceStatusResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(*picture.ServiceStatusResponse)
	return endpoints.ServiceStatusResponse{Code: response.GetCode(), Err: response.GetErr()}, nil
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
			encodeGRPCCreateRequest,
			decodeGRPCCreateResponse,
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

func encodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*endpoints.CreateRequest)
	newReq := &picture.CreateRequest{
		Text: req.Text,
		Fill: req.Fill,
		Pos:  picture.Position(picture.Position_value[string(req.Pos)]),
	}
	buf := new(bytes.Buffer)
	if req.Image != nil {
		png.Encode(buf, req.Image)
		newReq.Image = &picture.Image{Data: buf.Bytes(), Type: ".png"}
	}
	if req.Logo != nil {
		buf1 := new(bytes.Buffer)
		png.Encode(buf1, req.Logo)
		newReq.Logo = &picture.Image{Data: buf1.Bytes(), Type: ".png"}
	}
	return newReq, nil
}

func encodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	_ = grpcReq.(*endpoints.ServiceStatusRequest)
	return &picture.ServiceStatusRequest{}, nil
}

func decodeGRPCCreateResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*picture.CreateResponse)
	return &endpoints.CreateResponse{Image: getImageFromByte(resp.Image, ".png"), Err: resp.Err}, nil
}

func decodeGRPCServiceStatusResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(*picture.ServiceStatusResponse)
	return &endpoints.ServiceStatusResponse{Code: resp.GetCode(), Err: resp.GetErr()}, nil
}

func getImageFromByte(image []byte, encoding string) image.Image {
	switch encoding {
	case ".png":
		img, err := png.Decode(bytes.NewReader(image))
		if err != nil {
			return nil
		}
		return img
	case ".jpg":
		img, err := jpeg.Decode(bytes.NewReader(image))
		if err != nil {
			return nil
		}
		return img
	default:
		return nil
	}
}
