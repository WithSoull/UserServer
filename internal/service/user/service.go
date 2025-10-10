package user

import (
	"github.com/WithSoull/platform_common/pkg/client/db"
	"github.com/WithSoull/UserServer/internal/repository"
	"github.com/WithSoull/UserServer/internal/service"
)

type userService struct {
	repo     repository.UserRepository
	txManger db.TxManager
}

func NewService(
	repo repository.UserRepository,
	txManger db.TxManager,
) service.UserService {
	return &userService{
		repo:     repo,
		txManger: txManger,
	}
}
