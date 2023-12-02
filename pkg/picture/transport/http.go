package transport

import (
	"context"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"regexp"
	"watermark-service/internal"
	"watermark-service/internal/util"
	"watermark-service/pkg/picture"
	"watermark-service/pkg/picture/endpoints"

	httpkit "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	m := http.NewServeMux()

	m.Handle("/create", httpkit.NewServer(
		ep.CreateEndpoint,
		decodeHTTPCreateRequest,
		encodeCreateResponse,
		httpkit.ServerBefore(extractImages),
	))
	m.Handle("/healthz", httpkit.NewServer(
		ep.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
	))
	return m
}

func decodeHTTPServiceStatusRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req endpoints.ServiceStatusRequest
	return req, nil
}

func decodeHTTPCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateRequest
	val := ctx.Value(picture.ImageContextKey("image"))
	if val == nil {
		return nil, util.ErrInvalidArg
	}
	img, ok := val.(image.Image)
	if !ok {
		return nil, util.ErrInvalidArg
	}
	val = ctx.Value(picture.LogoContextKey("logo"))
	logo, ok := val.(image.Image)
	if !ok {
		req.Logo = nil
	}
	req.Image = img
	req.Logo = logo
	req.Fill = r.FormValue("fill") == "true"
	req.Text = r.FormValue("text")
	req.Pos = internal.PositionFromString(r.FormValue("pos"))
	return req, nil
}

func encodeCreateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if resp, ok := response.(endpoints.CreateResponse); ok {
		w.Header().Set("Content-Type", "image/png")
		return png.Encode(w, resp.Image)
	}
	return encodeResponse(ctx, w, response)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case util.ErrUnknownArg:
		w.WriteHeader(http.StatusNotFound)
	case util.ErrInvalidArg:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func extractImages(ctx context.Context, r *http.Request) (newCtx context.Context) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return ctx
	}
	img := getImageFromFile("image", r)
	newCtx = context.WithValue(ctx, picture.ImageContextKey("image"), img)
	logo := getImageFromFile("logo", r)
	newCtx = context.WithValue(newCtx, picture.LogoContextKey("logo"), logo)
	return
}

func getImageFromFile(name string, r *http.Request) image.Image {
	var image image.Image
	var regexp_result []string

	res, _ := regexp.Compile(`.[0-9a-z]+$`)
	file, header, err := r.FormFile(name)
	if err != nil {
		return nil
	}
	regexp_result = res.FindAllString(header.Filename, -1)
	if len(regexp_result) == 0 {
		return nil
	}
	switch regexp_result[0] {
	case ".png":
		image, err = png.Decode(file)
		if err != nil {
			image = nil
		}
	case ".jpg":
		image, err = jpeg.Decode(file)
		if err != nil {
			image = nil
		}
	default:
		image = nil
	}
	return image
}
