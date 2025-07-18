package user

import (
	"github.com/WithSoull/AuthService/internal/client/db"
	"github.com/WithSoull/AuthService/internal/repository"
)

type service struct {
	repo     repository.UserRepository
	txManger db.TxManager
}

func NewService(
	repo repository.UserRepository,
	txManger db.TxManager,
) *service {
	return &service{
		repo:     repo,
		txManger: txManger,
	}
}
