package transport

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"watermark-service/api/v1/protos/watermark"
	"watermark-service/internal"
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
	var Logo image.Image
	req := grpcReq.(*watermark.CreateRequest)
	img := req.GetImage()
	logo := req.GetLogo()
	if logo != nil {
		Logo = getImageFromByte(logo.Data, logo.Type)
	}
	pos := internal.PositionFromString(req.Pos.String())
	return endpoints.CreateRequest{Image: getImageFromByte(img.Data, img.Type), Logo: Logo, Text: req.Text, Fill: req.Fill, Pos: pos}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoints.ServiceStatusRequest{}, nil
}

func encodeGRPCCreateResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(endpoints.CreateResponse)
	buf := new(bytes.Buffer)
	png.Encode(buf, response.Image)
	return &watermark.CreateResponse{Image: buf.Bytes(), Err: response.Err}, nil
}

func encodeGRPCServiceStatusResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(*watermark.ServiceStatusResponse)
	return endpoints.ServiceStatusResponse{Code: response.GetCode(), Err: response.GetErr()}, nil
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
