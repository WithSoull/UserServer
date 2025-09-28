package repository

import (
	"context"

	"github.com/WithSoull/UserServer/internal/model"
)

type UserRepository interface {
	Create(context.Context, *model.UserInfo, string) (int64, error)
	Get(context.Context, int64) (*model.User, error)
	Update(context.Context, int64, *string, *string) error
	UpdatePassword(context.Context, int64, string) error
	LogPassword(context.Context, int64, string) error
	Delete(context.Context, int64) error
	GetUserCredentials(context.Context, string) (int64, string, error)
}
