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
)

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
