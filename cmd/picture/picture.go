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
	"github.com/kelseyhightower/envconfig"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

var (
	cfg config.PictureConfig
)

func main() {
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
	err := g.Run()
	zap.L().Info("exit", zap.Error(err))
}

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))

	var conf string
	flag.StringVar(&conf, "config", "", "config file")
	flag.Parse()
	if conf == "" {
		err := envconfig.Process("picture", &cfg)
		if err != nil {
			zap.L().Fatal("Setup failed", zap.String("config", "loading"), zap.Error(err))
		}
	} else {
		f, err := os.Open(conf)
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

}
