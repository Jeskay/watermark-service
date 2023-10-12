package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"watermark-service/internal/database"
	"watermark-service/internal/util"
	dbsvc "watermark-service/pkg/database"
	"watermark-service/pkg/database/endpoints"
	"watermark-service/pkg/database/transport"

	proto "watermark-service/api/v1/protos/db"

	grpckit "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultHTTPPort = "9091"
	defaultGRPCPort = "9092"
)

var (
	logger      log.Logger
	httpAddr    = net.JoinHostPort("localhost", util.EnvString("HTTP_PORT", defaultHTTPPort))
	grpcAddr    = net.JoinHostPort("localhost", util.EnvString("GRPC_PORT", defaultGRPCPort))
	authSvcAddr = net.JoinHostPort("localhost", util.EnvString("AUTH_SVC_PORT", "9022"))
)

func main() {

	orm, err := database.Init(database.DefaultHost, database.DefaultPort, database.DefaultDBUser, database.DefaultDatabase, database.DefaultPassword)
	if err != nil {
		logger.Log("FATAL: failed to load db with error ", err.Error())
	}

	var service dbsvc.Service
	{
		service = dbsvc.NewService(orm)
		service = dbsvc.AuthMiddleware(authSvcAddr)(service)
	}

	var (
		eps         = endpoints.NewEndpointSet(service)
		httpHandler = transport.NewHttpHandler(eps)
		grpcServer  = transport.NewGRPCServer(eps)
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
			proto.RegisterDatabaseServer(baseServer, grpcServer)
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
