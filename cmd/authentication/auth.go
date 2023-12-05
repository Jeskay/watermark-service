package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"watermark-service/config"
	"watermark-service/internal"
	authsvc "watermark-service/pkg/authentication"
	"watermark-service/pkg/authentication/endpoints"
	"watermark-service/pkg/authentication/transport"

	proto "watermark-service/api/v1/protos/auth"

	grpckit "github.com/go-kit/kit/transport/grpc"
	"github.com/kelseyhightower/envconfig"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

var (
	cfg config.AuthenticationConfig
)

func main() {
	var (
		grpcAddr = net.JoinHostPort(cfg.GRPCAddress.Host, cfg.GRPCAddress.Port)
		httpAddr = net.JoinHostPort(cfg.HTTPAddress.Host, cfg.HTTPAddress.Port)
	)
	connectionStr := internal.DatabaseConnectionStr{
		Host:     cfg.DbConnection.Host,
		Port:     cfg.DbConnection.Port,
		User:     cfg.DbConnection.User,
		Database: cfg.DbConnection.Database,
		Password: cfg.DbConnection.Password,
	}

	var service authsvc.Service
	{
		service = authsvc.NewService(connectionStr, cfg.SecretKey)
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
			zap.L().Fatal("transport", zap.String("HTTP", "during Listen"), zap.Error(err))
		}
		g.Add(func() error {
			zap.L().Info("transport", zap.String("HTTP", "Listener"), zap.String("address", httpAddr))
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			zap.L().Error("transport", zap.String("HTTP", "listener failed"), zap.Error(err))
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
			proto.RegisterAuthenticationServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(err error) {
			zap.L().Error("transport", zap.String("gRPC", "listener failed"), zap.Error(err))
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
	var conf, logf string
	flag.StringVar(&conf, "config", "", "config file")
	flag.StringVar(&logf, "log", "./var/logs", "log folder")
	flag.Parse()

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	filePath := logf + "/log-" + time.Now().String() + ".log"
	logFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)
	zap.ReplaceGlobals(zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)))

	if conf == "" {
		err := envconfig.Process("authentication", &cfg)
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
