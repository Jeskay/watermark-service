package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"watermark-service/config"
	"watermark-service/internal/database"
	dbsvc "watermark-service/pkg/database"
	"watermark-service/pkg/database/endpoints"
	"watermark-service/pkg/database/transport"

	proto "watermark-service/api/v1/protos/db"

	grpckit "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

var (
	cfg    config.DatabaseConfig
	logger log.Logger
)

func main() {
	f, err := os.Open("../../config/database_config.yaml")
	if err != nil {
		logger.Log("FATAL: failed to load config", err.Error())
	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Log("FATAL: failed to decode config file", err.Error())
	}
	f.Close()

	var (
		grpcAddr    = net.JoinHostPort(cfg.GRPCAddress.Host, cfg.GRPCAddress.Port)
		httpAddr    = net.JoinHostPort(cfg.HTTPAddress.Host, cfg.HTTPAddress.Port)
		authSvcAddr = net.JoinHostPort(cfg.Services.Auth.Host, cfg.Services.Auth.Port)
	)
	orm, err := database.Init(cfg.DbConnection.Host, cfg.DbConnection.Port, cfg.DbConnection.User, cfg.DbConnection.Database, cfg.DbConnection.Password)
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
