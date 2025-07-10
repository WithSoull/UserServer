package app

import (
	"context"
	"log"

	"github.com/WithSoull/AuthService/internal/closer"
	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/config/env"
	"github.com/WithSoull/AuthService/internal/repository"
	userRepository "github.com/WithSoull/AuthService/internal/repository/user"
	userService "github.com/WithSoull/AuthService/internal/service/user"
	userHandler "github.com/WithSoull/AuthService/internal/handler/user"
	"github.com/WithSoull/AuthService/internal/service"
	desc "github.com/WithSoull/AuthService/pkg/user/v1"
	"github.com/jackc/pgx/v5/pgxpool"
)

type serviceProvider struct {
	pgConfig config.PGCongif
	grpcConfig config.GRPCCongif

	pgPool *pgxpool.Pool

	userRepository repository.UserRepository
	userService service.UserService
	userHandler desc.UserV1Server
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGCongif {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCCongif {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) PGPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create connection pool: %s", err.Error())
		}

		if err := pool.Ping(ctx); err != nil {
			log.Fatalf("failed to connect to database: %v", err.Error())
		}

		closer.Add(func() error {
				pool.Close()
				return nil
			})

		s.pgPool = pool
	}

	return s.pgPool
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.PGPool(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserHandler(ctx context.Context) desc.UserV1Server {
	if s.userHandler == nil {
		s.userHandler = userHandler.NewHandler(s.UserService(ctx))
	}

	return s.userHandler
}
