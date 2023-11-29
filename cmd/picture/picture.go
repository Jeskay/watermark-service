package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	proto "watermark-service/api/v1/protos/picture"
	"watermark-service/config"
	"watermark-service/pkg/picture"
	"watermark-service/pkg/picture/endpoints"
	"watermark-service/pkg/picture/transport"

	grpckit "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

var (
	cfg    config.PictureConfig
	logger log.Logger
)

func main() {
	var build bool
	flag.BoolVar(&build, "built", false, "use context for built executable")
	flag.Parse()
	logger.Log(build)
	var confPath string
	if build {
		confPath = "./config/picture_config.yaml"
	} else {
		confPath = "../../config/picture_config.yaml"
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
		httpAddr = net.JoinHostPort(cfg.HTTPAddress.Host, cfg.HTTPAddress.Port)
		grpcAddr = net.JoinHostPort(cfg.GRPCAddress.Host, cfg.GRPCAddress.Port)
	)
	var service picture.Service
	{
		service = picture.NewService()
		service = picture.PictureMiddleware()(service)
	}

	var (
		eps         = endpoints.NewEndpointSet(service)
		httpHandler = transport.NewHTTPHandler(eps)
		grpcServer  = transport.NewGRPCServer(eps)
	)

	var g group.Group
	{
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
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
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", grpcAddr)
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(grpckit.Interceptor))
			proto.RegisterPictureServer(baseServer, grpcServer)
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
