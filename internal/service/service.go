package service

import (
	"context"

	"github.com/WithSoull/AuthService/internal/model"
)


type UserService interface {
	Create(ctx context.Context, info model.UserInfo, password, passwordConfirm string) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, id int64, name, email *string) error
	UpdatePassword(ctx context.Context, id int64, password, confirm_password string) (error)
	Delete(ctx context.Context, id int64) error
}
