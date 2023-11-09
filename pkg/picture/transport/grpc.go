package transport

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"watermark-service/api/v1/protos/picture"
	"watermark-service/internal"
	"watermark-service/pkg/picture/endpoints"

	grpckit "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	create        grpckit.Handler
	serviceStatus grpckit.Handler
	picture.UnimplementedPictureServer
}

func NewGRPCServer(ep endpoints.Set) picture.PictureServer {
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
	return endpoints.CreateRequest{Image: getImageFromByte(img.Data, img.Type), Logo: Logo, Text: req.Text, Fill: req.Fill, Pos: pos}, nil
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
