package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	auth "watermark-service/internal/authentication"
	"watermark-service/internal/util"
	authsvc "watermark-service/pkg/authentication"
	"watermark-service/pkg/authentication/endpoints"
	"watermark-service/pkg/authentication/transport"

	proto "watermark-service/api/v1/protos/auth"

	grpckit "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultHTTPPort = "9021"
	defaultGRPCPort = "9022"
)

var (
	logger   log.Logger
	grpcAddr = net.JoinHostPort("localhost", util.EnvString("GRPC_PORT", defaultGRPCPort))
	httpAddr = net.JoinHostPort("localhost", util.EnvString("HTTP_PORT", defaultHTTPPort))
)

func main() {
	orm, err := auth.Init(auth.DefaultHost, auth.DefaultPort, auth.DefaultDBUser, auth.DefaultDatabase, auth.DefaultPassword)
	if err != nil {
		logger.Log("FATAL: failed to load db with error ", err.Error())
	}

	var (
		service     = authsvc.NewService(orm, "SECRET_TOKEN")
		eps         = endpoints.NewEndpointSet(service)
		grpcServer  = transport.NewGRPCServer(eps)
		httpHandler = transport.NewHTTPHandler(eps)
	)

	var g group.Group
	{
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "error", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "error", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", grpcAddr)
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(grpckit.Interceptor))
			reflection.Register(baseServer)
			proto.RegisterAuthenticationServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
