package app

import (
	"context"
	"log"

	"github.com/WithSoull/AuthService/internal/client/db"
	"github.com/WithSoull/AuthService/internal/client/db/pg"
	"github.com/WithSoull/AuthService/internal/client/db/transaction"
	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/config/env"
	userHandler "github.com/WithSoull/AuthService/internal/handler/user"
	"github.com/WithSoull/AuthService/internal/repository"
	userRepository "github.com/WithSoull/AuthService/internal/repository/user"
	"github.com/WithSoull/AuthService/internal/service"
	userService "github.com/WithSoull/AuthService/internal/service/user"
	desc "github.com/WithSoull/AuthService/pkg/user/v1"
	"github.com/WithSoull/platform_common/pkg/closer"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	httpConfig config.HTTPConfig

	pgClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository
	userService    service.UserService
	userHandler    desc.UserV1Server
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) PGClient(ctx context.Context) db.Client {
	if s.pgClient == nil {
		client, err := pg.NewPGClient(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create connection pool: %s", err.Error())
		}

		if err := client.DB().Ping(ctx); err != nil {
			log.Fatalf("failed to connect to database: %v", err.Error())
		}

		closer.Add(func() error {
			client.Close()
			return nil
		})

		s.pgClient = client
	}

	return s.pgClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.PGClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.PGClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx), s.TxManager(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserHandler(ctx context.Context) desc.UserV1Server {
	if s.userHandler == nil {
		s.userHandler = userHandler.NewHandler(s.UserService(ctx))
	}

	return s.userHandler
}
