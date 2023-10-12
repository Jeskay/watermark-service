package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"watermark-service/pkg/authentication/endpoints"

	"watermark-service/internal/util"

	httpkit "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	m := http.NewServeMux()

	m.Handle("/healthz", httpkit.NewServer(
		ep.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
	))

	m.Handle("/login", httpkit.NewServer(
		ep.LoginEndpoint,
		decodeHTTPLoginRequest,
		encodeResponse,
	))

	m.Handle("/register", httpkit.NewServer(
		ep.RegisterEndpoint,
		decodeHTTPRegisterRequest,
		encodeResponse,
	))

	m.Handle("/validate", httpkit.NewServer(
		ep.ValidateEndpoint,
		decodeHTTPValidateRequest,
		encodeResponse,
	))

	return m
}

func decodeHTTPLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPRegisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPValidateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ValidateRequest
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
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	switch err {
	case util.ErrInvalidArg:
		w.WriteHeader(http.StatusNotFound)
	case util.ErrInvalidArg:
		w.WriteHeader(http.StatusBadRequest)
	default:
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
	}
}
