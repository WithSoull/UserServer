package service

import (
	"context"

	"github.com/WithSoull/UserServer/internal/model"
)

type UserService interface {
	Create(context.Context, model.UserInfo, string, string) (int64, error)
	Get(context.Context, int64) (*model.User, error)
	Update(context.Context, int64, *string, *string) error
	UpdatePassword(context.Context, int64, string, string) error
	Delete(context.Context, int64) error
	ValidateCredentials(context.Context, string, string) (bool, int64)
}
