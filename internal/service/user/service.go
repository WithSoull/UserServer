package user

import (
	"github.com/WithSoull/UserServer/internal/repository"
	"github.com/WithSoull/UserServer/internal/service"
	"github.com/WithSoull/platform_common/pkg/client/db"
)

type userService struct {
	repo     repository.UserRepository
	txManger db.TxManager

	userProducerService service.UserProducerService
}

func NewService(
	repo repository.UserRepository,
	txManger db.TxManager,
	userProducerService service.UserProducerService,
) service.UserService {
	return &userService{
		repo:                repo,
		txManger:            txManger,
		userProducerService: userProducerService,
	}
}
