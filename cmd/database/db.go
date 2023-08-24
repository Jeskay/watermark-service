package main

import (
	"net"
	"os"

	"github.com/go-kit/log"
)

const (
	defaultHTTPPort = "8081"
	defaultGRPCPort = "8082"
)

var (
	logger log.Logger
	httpAddr = net.JoinHostPort("localhost", envString("HTTP_PORT", defaultHTTPPort))
	grpcAddr = net.JoinHostPort("localhost", envString("GRPC_PORT", defaultGRPCPort))
)

func envString(env, default_value string) string {
	e := os.Getenv(env)
	if e == "" {
		return default_value
	}
	return e
}

func init() {
	logger  = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	db, err := database.
}