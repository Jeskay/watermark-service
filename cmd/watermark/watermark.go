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
	"watermark-service/internal"
	watermarksvc "watermark-service/pkg/watermark"
	"watermark-service/pkg/watermark/endpoints"
	"watermark-service/pkg/watermark/transport"

	proto "watermark-service/api/v1/protos/watermark"

	grpckit "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

var (
	cfg config.WatermarkConfig
)

func main() {
	var (
		grpcAddr       = net.JoinHostPort(cfg.GRPCAddress.Host, cfg.GRPCAddress.Port)
		httpAddr       = net.JoinHostPort(cfg.HTTPAddress.Host, cfg.HTTPAddress.Port)
		authSvcAddr    = net.JoinHostPort(cfg.Services.Auth.Host, cfg.Services.Auth.Port)
		pictureSvcAddr = net.JoinHostPort(cfg.Services.Picture.Host, cfg.Services.Picture.Port)
	)
	connectionStr := internal.DatabaseConnectionStr{
		Host:     cfg.DbConnection.Host,
		Port:     cfg.DbConnection.Port,
		User:     cfg.DbConnection.User,
		Database: cfg.DbConnection.Database,
		Password: cfg.DbConnection.Password,
	}

	var service watermarksvc.Service
	{
		service = watermarksvc.NewService(connectionStr, pictureSvcAddr, cfg.Cloudinary.Cloud, cfg.Cloudinary.Api, cfg.Cloudinary.Secret)
		service = watermarksvc.AuthMiddleware(authSvcAddr)(service)
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
			zap.L().Fatal("transport", zap.String("HTTP", "during Listen"), zap.Error(err))
		}
		g.Add(func() error {
			zap.L().Info("transport", zap.String("HTTP", "Listener"), zap.String("address", httpAddr))
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			zap.L().Fatal("transport", zap.String("gRPC", "during Listen"), zap.Error(err))
		}
		g.Add(func() error {
			zap.L().Info("transport", zap.String("gRPC", "Listener"), zap.String("address", grpcAddr))
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(grpckit.Interceptor))
			reflection.Register(baseServer)
			proto.RegisterWatermarkServer(baseServer, grpcServer)
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
	err := g.Run()
	zap.L().Info("exit", zap.Error(err))
}

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))

	var build bool
	flag.BoolVar(&build, "built", false, "use context for built executable")
	flag.Parse()

	var confPath string
	if build {
		confPath = "./config/watermark_config.yaml"
	} else {
		confPath = "../../config/watermark_config.yaml"
	}
	f, err := os.Open(confPath)
	if err != nil {
		zap.L().Fatal("Setup failed", zap.String("config", "loading"), zap.Error(err))
	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		zap.L().Fatal("Setup failed", zap.String("config", "decoding"), zap.Error(err))
	}
	f.Close()
}
