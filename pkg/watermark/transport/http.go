package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"watermark-service/internal/util"
	"watermark-service/pkg/watermark/endpoints"

	httpkit "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/log"
)

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	m := http.NewServeMux()

	m.Handle("/healthz", httpkit.NewServer(
		ep.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
	))
	m.Handle("/status", httpkit.NewServer(
		ep.StatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
	))
	m.Handle("/addDocument", httpkit.NewServer(
		ep.AddDocumentEndpoint,
		decodeHTTPAddDocumentRequest,
		encodeResponse,
	))
	m.Handle("/get", httpkit.NewServer(
		ep.GetEndpoint,
		decodeHTTPGetRequest,
		encodeResponse,
	))
	m.Handle("/watermark", httpkit.NewServer(
		ep.WatermarkEndpoint,
		decodeHTTPWatermarkRequest,
		encodeResponse,
	))
	return m
}

func decodeHTTPGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetRequest
	if r.ContentLength == 0 {
		logger.Log("Get request with empty body")
		return req, nil
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.StatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPWatermarkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.WatermarkRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPAddDocumentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.AddDocumentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPServiceStatusRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req endpoints.ServiceStatusRequest
	return req, nil
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

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
