package app

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"sync"
	"syscall"
	"time"

	"github.com/WithSoull/UserServer/internal/config"
	desc "github.com/WithSoull/UserServer/pkg/user/v1"
	"github.com/WithSoull/platform_common/pkg/closer"
	"github.com/WithSoull/platform_common/pkg/logger"
	validationInterceptor "github.com/WithSoull/platform_common/pkg/middleware/validation"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	shutdownTimeout = 10 * time.Second
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initCloser,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(config.AppConfig().Logger.LogLevel(), config.AppConfig().Logger.AsJSON())
}

func (a *App) initCloser(_ context.Context) error {
	closer.Configure(logger.Logger(), shutdownTimeout, syscall.SIGINT, syscall.SIGTERM)
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	creds, err := credentials.NewServerTLSFromFile("service.pem", "service.key")
	if err != nil {
		log.Fatalf("failed to load TLS keys: %v", err)
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				validationInterceptor.ErrorCodesInterceptor,
			),
		),
	)

	closer.AddNamed("GRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserHandler(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()
	creds, err := credentials.NewServerTLSFromFile("service.pem", "service.key")
	if err != nil {
		log.Fatalf("failed to load TLS keys: %v", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	err = desc.RegisterUserV1HandlerFromEndpoint(ctx, mux, config.AppConfig().GRPC.Address(), opts)
	if err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    config.AppConfig().HTTP.Address(),
		Handler: mux,
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", config.AppConfig().GRPC.Address())
	if err != nil {
		return err
	}

	logger.Info(context.Background(), "GRPC server listening", zap.String("address", config.AppConfig().GRPC.Address()))

	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	logger.Info(context.Background(), "GRPC server stopped gracefully")
	return nil
}

func (a *App) runHTTPServer() error {
	logger.Info(context.Background(), "HTTP server is starting listening and serving", zap.String("address", config.AppConfig().HTTP.Address()))
	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	logger.Info(context.Background(), "HTTP server stopped gracefully")
	return nil
}

func (a *App) Run() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			logger.Error(context.Background(), "fault grpc server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(context.Background(), "fault http server", zap.Error(err))
		}
	}()

	wg.Wait()

	return nil
}
