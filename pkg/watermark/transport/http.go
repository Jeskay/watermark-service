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
	"watermark-service/pkg/watermark/endpoints"

	httpkit "github.com/go-kit/kit/transport/http"
)

func NewHttpHandler(ep endpoints.Set) http.Handler {
	m := http.NewServeMux()

	opts := []httpkit.ServerOption{
		httpkit.ServerBefore(injectContext),
	}

	m.Handle("/healthz", httpkit.NewServer(
		ep.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
	))
	m.Handle("/add", httpkit.NewServer(
		ep.AddEndpoint,
		decodeHTTPAddRequest,
		encodeResponse,
		httpkit.ServerBefore(injectContext, extractImages),
	))
	m.Handle("/get", httpkit.NewServer(
		ep.GetEndpoint,
		decodeHTTPGetRequest,
		encodeResponse,
		opts...,
	))
	m.Handle("/remove", httpkit.NewServer(
		ep.RemoveEndpoint,
		decodeHTTPRemoveRequest,
		encodeResponse,
		opts...,
	))

	return m
}

func decodeHTTPGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetRequest
	if r.ContentLength == 0 {
		return req, nil
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPRemoveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.RemoveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPAddRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.AddRequest
	val := ctx.Value("image")
	if val == nil {
		return nil, util.ErrInvalidArg
	}
	img, ok := val.(image.Image)
	if !ok {
		return nil, util.ErrInvalidArg
	}
	val = ctx.Value("logo")
	logo, ok := val.(image.Image)
	if !ok {
		return nil, util.ErrInvalidArg
	}
	req.Image = img
	req.Logo = logo
	req.Fill = r.FormValue("fill") == "true"
	req.Text = r.FormValue("text")
	req.Pos = internal.PositionFromString(r.FormValue("pos"))

	return req, nil
}

func decodeHTTPServiceStatusRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req endpoints.ServiceStatusRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
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
	newCtx = context.WithValue(ctx, "image", img)
	logo := getImageFromFile("logo", r)
	newCtx = context.WithValue(newCtx, "logo", logo)
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

func injectContext(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, "token", r.Header.Get("Token"))
}
