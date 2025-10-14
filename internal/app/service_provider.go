package app

import (
	"context"

	"github.com/WithSoull/UserServer/internal/config"
	userHandler "github.com/WithSoull/UserServer/internal/handler/user"
	"github.com/WithSoull/UserServer/internal/repository"
	userRepository "github.com/WithSoull/UserServer/internal/repository/user"
	"github.com/WithSoull/UserServer/internal/service"
	userService "github.com/WithSoull/UserServer/internal/service/user"
	desc "github.com/WithSoull/UserServer/pkg/user/v1"
	"github.com/WithSoull/platform_common/pkg/client/db"
	"github.com/WithSoull/platform_common/pkg/client/db/pg"
	"github.com/WithSoull/platform_common/pkg/client/db/transaction"
	"github.com/WithSoull/platform_common/pkg/closer"
	"github.com/WithSoull/platform_common/pkg/logger"
)

type serviceProvider struct {
	pgClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository
	userService    service.UserService
	userHandler    desc.UserV1Server
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGClient(ctx context.Context) db.Client {
	if s.pgClient == nil {
		client, err := pg.NewPGClient(ctx, config.AppConfig().PG.DSN(), logger.Logger())
		if err != nil {
			panic(err)
		}

		if err := client.DB().Ping(ctx); err != nil {
			panic(err)
		}

		closer.AddNamed("PGClient", func(ctx context.Context) error {
			return client.Close()
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
