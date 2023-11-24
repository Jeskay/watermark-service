package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"watermark-service/config"
	auth "watermark-service/internal/authentication"
	authsvc "watermark-service/pkg/authentication"
	"watermark-service/pkg/authentication/endpoints"
	"watermark-service/pkg/authentication/transport"

	proto "watermark-service/api/v1/protos/auth"

	grpckit "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

var (
	cfg    config.AuthenticationConfig
	logger log.Logger
)

func main() {
	build := flag.Bool("build", false, "use context for built executable")
	flag.Parse()
	var confPath string
	if *build {
		confPath = "/config/auth_config.yaml"
	} else {
		confPath = "../../config/auth_config.yaml"
	}
	f, err := os.Open(confPath)
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
		grpcAddr = net.JoinHostPort(cfg.GRPCAddress.Host, cfg.GRPCAddress.Port)
		httpAddr = net.JoinHostPort(cfg.HTTPAddress.Host, cfg.HTTPAddress.Port)
	)
	orm, err := auth.Init(cfg.DbConnection.Host, cfg.DbConnection.Port, cfg.DbConnection.User, cfg.DbConnection.Database, cfg.DbConnection.Password)
	if err != nil {
		logger.Log("FATAL: failed to load db with error ", err.Error())
	}

	var service authsvc.Service
	{
		service = authsvc.NewService(orm, cfg.SecretKey)
		service = authsvc.AuthMiddleware()(service)
	}

	var (
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
