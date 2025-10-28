package app

import (
	"context"
	"fmt"
	"os"

	"github.com/IBM/sarama"
	"github.com/WithSoull/UserServer/internal/config"
	userHandler "github.com/WithSoull/UserServer/internal/handler/user"
	"github.com/WithSoull/UserServer/internal/repository"
	userRepository "github.com/WithSoull/UserServer/internal/repository/user"
	"github.com/WithSoull/UserServer/internal/service"
	userProducerService "github.com/WithSoull/UserServer/internal/service/producer/user"
	userService "github.com/WithSoull/UserServer/internal/service/user"
	desc "github.com/WithSoull/UserServer/pkg/user/v1"
	"github.com/WithSoull/platform_common/pkg/client/db"
	"github.com/WithSoull/platform_common/pkg/client/db/pg"
	"github.com/WithSoull/platform_common/pkg/client/db/transaction"
	"github.com/WithSoull/platform_common/pkg/closer"
	"github.com/WithSoull/platform_common/pkg/kafka"
	kafkaProducer "github.com/WithSoull/platform_common/pkg/kafka/producer"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

type serviceProvider struct {
	pgClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository

	userService         service.UserService
	userProducerService service.UserProducerService

	userHandler desc.UserV1Server

	syncProducer        sarama.SyncProducer
	userCreatedProducer kafka.Producer
	userDeletedProducer kafka.Producer
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGClient(ctx context.Context) db.Client {
	if s.pgClient == nil {
		client, err := pg.NewPGClient(ctx, logger.Logger(), config.AppConfig().PG)
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
		s.userService = userService.NewService(s.UserRepository(ctx), s.TxManager(ctx), s.UserProducerService(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserHandler(ctx context.Context) desc.UserV1Server {
	if s.userHandler == nil {
		s.userHandler = userHandler.NewHandler(s.UserService(ctx))
	}

	return s.userHandler
}

func (s *serviceProvider) SyncProducer(ctx context.Context) sarama.SyncProducer {
	if s.syncProducer == nil {
		logger.Info(ctx,
			"Creating Kafka SyncProducer",
			zap.Strings("brokers", config.AppConfig().Kafka.Brokers()),
			zap.String("raw_env", os.Getenv("KAFKA_BROKERS")), // добавь это
		)

		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().Sarama.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %v\n", err))
		}

		closer.AddNamed(
			"KafkaSyncProducer",
			func(ctx context.Context) error {
				return p.Close()
			},
		)
		s.syncProducer = p
	}

	return s.syncProducer
}

func (s *serviceProvider) UserCreatedProducer(ctx context.Context) kafka.Producer {
	if s.userCreatedProducer == nil {
		s.userCreatedProducer = kafkaProducer.NewProducer(
			s.SyncProducer(ctx),
			config.AppConfig().UserCreatedProducer.Topic(),
			logger.Logger(),
		)
	}

	return s.userCreatedProducer
}

func (s *serviceProvider) UserDeletedProducer(ctx context.Context) kafka.Producer {
	if s.userDeletedProducer == nil {
		s.userDeletedProducer = kafkaProducer.NewProducer(
			s.SyncProducer(ctx),
			config.AppConfig().UserDeletedProducer.Topic(),
			logger.Logger(),
		)
	}

	return s.userDeletedProducer
}

func (s *serviceProvider) UserProducerService(ctx context.Context) service.UserProducerService {
	if s.userProducerService == nil {
		s.userProducerService = userProducerService.NewService(
			s.UserCreatedProducer(ctx),
			s.UserDeletedProducer(ctx),
		)
	}

	return s.userProducerService
}
